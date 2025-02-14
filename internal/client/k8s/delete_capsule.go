package k8s

import (
	"context"
	"fmt"

	"github.com/rigdev/rig/pkg/auth"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeleteCapsule implements cluster.Gateway.
func (c *Client) DeleteCapsule(ctx context.Context, capsuleName string) error {
	projectID, err := auth.GetProjectID(ctx)
	if err != nil {
		return err
	}
	ns := projectID.String()

	if err := c.deleteProxyEnvSecret(ctx, capsuleName, ns); err != nil {
		return err
	}
	if err := c.deleteLoadBalancer(ctx, capsuleName, ns); err != nil {
		return err
	}
	if err := c.deleteIngress(ctx, capsuleName, ns); err != nil {
		return err
	}
	if err := c.deleteService(ctx, capsuleName, ns); err != nil {
		return err
	}
	if err := c.deleteEnvSecret(ctx, capsuleName, ns); err != nil {
		return err
	}
	if err := c.deleteDeployment(ctx, capsuleName, ns); err != nil {
		return err
	}

	return nil
}

func (c *Client) deletePullSecret(ctx context.Context, namespace string) error {
	err := c.cs.CoreV1().Secrets(namespace).
		Delete(ctx, fmt.Sprintf("%s-pull", namespace), metav1.DeleteOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("could not delete pull Secret: %w", err)
	}
	return nil

}

func (c *Client) deleteProxyEnvSecret(ctx context.Context, capsuleName, namespace string) error {
	err := c.cs.CoreV1().Secrets(namespace).
		Delete(ctx, fmt.Sprintf("%s-proxy", capsuleName), metav1.DeleteOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("could not delete proxy env Secret: %w", err)
	}
	return nil
}

func (c *Client) deleteLoadBalancer(ctx context.Context, capsuleName, namespace string) error {
	err := c.cs.CoreV1().Services(namespace).
		Delete(ctx, fmt.Sprintf("%s-lb", capsuleName), metav1.DeleteOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("could not delete load balancer Service: %w", err)
	}
	return nil
}

func (c *Client) deleteIngress(ctx context.Context, capsuleName, ns string) error {
	err := c.cs.NetworkingV1().
		Ingresses(ns).
		Delete(ctx, capsuleName, metav1.DeleteOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("could not delete Ingress: %w", err)
	}
	return nil
}

func (c *Client) deleteService(ctx context.Context, capsuleName, ns string) error {
	err := c.cs.CoreV1().
		Services(ns).
		Delete(ctx, capsuleName, metav1.DeleteOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("could not delete Service: %w", err)
	}
	return nil
}

func (c *Client) deleteEnvSecret(ctx context.Context, capsuleName, ns string) error {
	err := c.cs.CoreV1().
		Secrets(ns).
		Delete(ctx, capsuleName, metav1.DeleteOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("could not delete Secret: %w", err)
	}
	return nil
}

func (c *Client) deleteConfigMap(ctx context.Context, capsuleName, ns string) error {
	err := c.cs.CoreV1().
		ConfigMaps(ns).
		Delete(ctx, capsuleName, metav1.DeleteOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("could not delete ConfigMap: %w", err)
	}
	return nil
}

func (c *Client) deleteDeployment(ctx context.Context, capsuleName, ns string) error {
	err := c.cs.AppsV1().
		Deployments(ns).
		Delete(ctx, capsuleName, metav1.DeleteOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("could not delete Deployment: %w", err)
	}
	return nil
}
