package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	createName     string
	createImage    string
	createReplicas int32
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Deployment",
	Run: func(cmd *cobra.Command, args []string) {
		clientset := mustGetClientSet()
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: createName,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &createReplicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": createName},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": createName},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  createName,
							Image: createImage,
						}},
					},
				},
			},
		}

		_, err := clientset.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create Deployment")
		}
		log.Info().Msgf("✅ Deployment %q created", createName)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createName, "name", "", "Name of the Deployment")
	createCmd.Flags().StringVar(&createImage, "image", "nginx:latest", "Container image")
	createCmd.Flags().Int32Var(&createReplicas, "replicas", 1, "Number of replicas")
	if err := createCmd.MarkFlagRequired("name"); err != nil {
		log.Warn().Err(err).Msg("Failed to mark 'name' as required")
	}
}
