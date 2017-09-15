package action_test

import (
	. "github.com/cloudfoundry-incubator/bosh-load-tests/action"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"fmt"
)

var _ = Describe("DeployWithStatic", func() {

	It("should get next ip correctly", func() {
		ips := make(map[string]struct{})
		for j := 0; j < 300; j++ {
			withStatic := NewDeployWithStatic(DirectorInfo{}, j, fmt.Sprintf("deploymentName_%d", j), nil, nil, nil, false)
			for i := 0; i < 10; i ++ {
				ip := withStatic.GetNextIP(i)
				if _, found := ips[ip]; found {
					Expect(found).To(BeFalse(), fmt.Sprintf("duplicate found for deployment %d, ip address %+v", j, ip))
				} else {
					ips[ip] = struct{}{}
				}
			}
		}
		Expect(ips).To(HaveLen(10 * 300))

	})
})
