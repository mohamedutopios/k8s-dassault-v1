package controller

import (
    "context"
    "encoding/json"
    "time"

    "github.com/go-logr/logr"
    "github.com/go-redis/redis/v8"
    "k8s.io/apimachinery/pkg/api/errors"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"

	appsv1alpha1 "vote-campaign-operator/api/v1alpha1"
)

var ctx = context.Background()

// VoteCampaignReconciler reconciles a VoteCampaign object
type VoteCampaignReconciler struct {
    client.Client
    Log    logr.Logger
    Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps.example.com,resources=votecampaigns,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.example.com,resources=votecampaigns/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.example.com,resources=votecampaigns/finalizers,verbs=update

func (r *VoteCampaignReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    r.Log.WithValues("votecampaign", req.NamespacedName).Info("Reconciling VoteCampaign")

    // Fetch the VoteCampaign instance
    campaign := &appsv1alpha1.VoteCampaign{}
    err := r.Get(ctx, req.NamespacedName, campaign)
    if err != nil {
        if errors.IsNotFound(err) {
            // Request object not found, could have been deleted after reconcile request.
            // Return and don't requeue
            r.Log.Info("VoteCampaign resource not found. Ignoring since object must be deleted")
            return ctrl.Result{}, nil
        }
        // Error reading the object - requeue the request.
        r.Log.Error(err, "Failed to get VoteCampaign")
        return ctrl.Result{}, err
    }

    // Check if the campaign is active
    now := time.Now()
    startTime, _ := time.Parse(time.RFC3339, campaign.Spec.StartTime)
    endTime, _ := time.Parse(time.RFC3339, campaign.Spec.EndTime)

    r.Log.Info("Current time", "now", now)
    r.Log.Info("Campaign start time", "startTime", startTime)
    r.Log.Info("Campaign end time", "endTime", endTime)

    if now.After(startTime) && now.Before(endTime) {
        campaign.Status.Active = true
        r.Log.Info("Campaign is active, retrieving votes")
        r.retrieveVotes(campaign)
    } else {
        campaign.Status.Active = false
        r.Log.Info("Campaign is not active")
    }

    // Update the status
    r.Log.Info("Updating VoteCampaign status")
    err = r.Status().Update(ctx, campaign)
    if err != nil {
        r.Log.Error(err, "Failed to update VoteCampaign status")
        return ctrl.Result{}, err
    }

    return ctrl.Result{RequeueAfter: time.Minute}, nil
}

func (r *VoteCampaignReconciler) retrieveVotes(campaign *appsv1alpha1.VoteCampaign) {
    r.Log.Info("Connecting to Redis")
    rdb := redis.NewClient(&redis.Options{
        Addr: "redis-service:6379",
    })

    votes := make(map[string]int)
    for _, option := range campaign.Spec.Options {
        votes[option.Name] = 0
    }

    r.Log.Info("Retrieving votes from Redis")
    voteData, err := rdb.LRange(ctx, "votes", 0, -1).Result()
    if err != nil {
        r.Log.Error(err, "Failed to retrieve votes from Redis")
        return
    }

    r.Log.Info("Votes retrieved from Redis", "voteData", voteData)

    for _, voteJSON := range voteData {
        var vote map[string]string
        err := json.Unmarshal([]byte(voteJSON), &vote)
        if err != nil {
            r.Log.Error(err, "Failed to unmarshal vote")
            continue
        }
        r.Log.Info("Vote unmarshaled", "vote", vote)
        if voteOption, ok := vote["vote"]; ok {
            if _, exists := votes[voteOption]; exists {
                votes[voteOption]++
                r.Log.Info("Vote counted", "voteOption", voteOption, "count", votes[voteOption])
            }
        }
    }

    campaign.Status.Votes = votes
    r.Log.Info("Votes retrieved and status updated", "votes", votes)
}

// SetupWithManager sets up the controller with the Manager.
func (r *VoteCampaignReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&appsv1alpha1.VoteCampaign{}).
        Complete(r)
}
