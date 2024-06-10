package controller

import (
    "path/filepath"
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "sigs.k8s.io/controller-runtime/pkg/envtest/printer"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/envtest"
    "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
    RegisterFailHandler(Fail)

    RunSpecsWithDefaultAndCustomReporters(t,
        "Controller Suite",
        []Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
    ctrl.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

    By("bootstrapping test environment")
    testEnv = &envtest.Environment{
        CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
    }

    var err error
    cfg, err := testEnv.Start()
    Expect(err).ToNot(HaveOccurred())
    Expect(cfg).ToNot(BeNil())

    k8sClient, err := client.New(cfg, client.Options{Scheme: scheme})
    Expect(err).ToNot(HaveOccurred())
    Expect(k8sClient).ToNot(BeNil())

    close(done)
}, 60)

var _ = AfterSuite(func() {
    By("tearing down the test environment")
    err := testEnv.Stop()
    Expect(err).ToNot(HaveOccurred())
})
