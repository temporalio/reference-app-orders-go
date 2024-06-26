#!/bin/sh

CLUSTER_NAME=temporal-oms

echo "Creating cluster: '$CLUSTER_NAME'..."

eksctl create cluster --name $CLUSTER_NAME

echo "Configuring ebs-csi-driver for volume storage..."

eksctl utils associate-iam-oidc-provider --region=us-west-2 --cluster=$CLUSTER_NAME --approve

eksctl create iamserviceaccount \
    --name ebs-csi-controller-sa \
    --namespace kube-system \
    --cluster $CLUSTER_NAME \
    --role-name AmazonEKS_EBS_CSI_DriverRole \
    --role-only \
    --attach-policy-arn arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy \
    --approve

eksctl create addon \
    --name aws-ebs-csi-driver \
    --cluster $CLUSTER_NAME \
    --service-account-role-arn arn:aws:iam::$(aws sts get-caller-identity --query Account --output text):role/AmazonEKS_EBS_CSI_DriverRole \
    --force

kubectl patch storageclass gp2 -p '{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"true"}}}'

echo "Done."