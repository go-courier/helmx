package spec

import (
	"testing"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
)

func TestRoleRule(t *testing.T) {
	t.Run("parse & string full", func(t *testing.T) {
		r, _ := ParseRoleRule("apps,extensions.deployments=a,b,c,d#update,get")

		gomega.NewWithT(t).Expect(r.ApiGroups).To(gomega.Equal([]string{"apps", "extensions"}))
		gomega.NewWithT(t).Expect(r.Resources).To(gomega.Equal([]string{"deployments"}))
		gomega.NewWithT(t).Expect(r.Verbs).To(gomega.Equal([]string{"update", "get"}))
		gomega.NewWithT(t).Expect(r.ResourceNames).To(gomega.Equal([]string{"a", "b", "c", "d"}))

		require.Equal(t, "apps,extensions.deployments=a,b,c,d#update,get", r.String())
	})

	t.Run("parse & string simple", func(t *testing.T) {
		r, _ := ParseRoleRule("secrets#get,update")

		gomega.NewWithT(t).Expect(r.ApiGroups).To(gomega.Equal([]string{""}))
		gomega.NewWithT(t).Expect(r.Resources).To(gomega.Equal([]string{"secrets"}))
		gomega.NewWithT(t).Expect(r.Verbs).To(gomega.Equal([]string{"get", "update"}))
		gomega.NewWithT(t).Expect(r.ResourceNames).To(gomega.BeNil())

		require.Equal(t, "secrets#get,update", r.String())
	})
}
