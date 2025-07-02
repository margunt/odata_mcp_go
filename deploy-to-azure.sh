#!/bin/bash

# Azure Container Registry and Container Instance Deployment Script
# This script builds and deploys the OData MCP bridge to Azure

set -e

# Configuration
RESOURCE_GROUP_NAME="rg-odata-mcp"
LOCATION="swedencentral"
ACR_NAME="odatamcpregistry$(date +%s)"  # Unique name
CONTAINER_NAME="odata-mcp-bridge"
IMAGE_NAME="odata-mcp"
TAG="latest"

echo "🚀 Starting Azure deployment for OData MCP Bridge..."

# Check if Azure CLI is installed
if ! command -v az &> /dev/null; then
    echo "❌ Azure CLI is not installed. Please install it first."
    echo "   Visit: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli"
    exit 1
fi

# Login check
echo "🔐 Checking Azure login status..."
if ! az account show &> /dev/null; then
    echo "🔑 Please login to Azure..."
    az login
fi

# Create Resource Group
echo "📦 Creating resource group: $RESOURCE_GROUP_NAME"
az group create \
    --name "$RESOURCE_GROUP_NAME" \
    --location "$LOCATION" \
    --output table

# Create Azure Container Registry
echo "🏗️  Creating Azure Container Registry: $ACR_NAME"
az acr create \
    --resource-group "$RESOURCE_GROUP_NAME" \
    --name "$ACR_NAME" \
    --sku Basic \
    --admin-enabled true \
    --output table

# Get ACR login server
ACR_LOGIN_SERVER=$(az acr show --name "$ACR_NAME" --resource-group "$RESOURCE_GROUP_NAME" --query "loginServer" --output tsv)
echo "📝 ACR Login Server: $ACR_LOGIN_SERVER"

# Build and push Docker image
echo "🔨 Building Docker image..."
docker build -t "$IMAGE_NAME:$TAG" .

echo "🏷️  Tagging image for ACR..."
docker tag "$IMAGE_NAME:$TAG" "$ACR_LOGIN_SERVER/$IMAGE_NAME:$TAG"

echo "🔐 Logging into ACR..."
az acr login --name "$ACR_NAME"

echo "📤 Pushing image to ACR..."
docker push "$ACR_LOGIN_SERVER/$IMAGE_NAME:$TAG"

# Update the parameter file with the correct image
echo "📝 Updating deployment parameters..."
sed -i "s|your-registry.azurecr.io/odata-mcp:latest|$ACR_LOGIN_SERVER/$IMAGE_NAME:$TAG|g" azure-parameters.json

# Get ACR credentials for the container instance
ACR_USERNAME=$(az acr credential show --name "$ACR_NAME" --query "username" --output tsv)
ACR_PASSWORD=$(az acr credential show --name "$ACR_NAME" --query "passwords[0].value" --output tsv)

# Deploy Container Instance
echo "🚀 Deploying Container Instance..."
DEPLOYMENT_OUTPUT=$(az deployment group create \
    --resource-group "$RESOURCE_GROUP_NAME" \
    --template-file azure-container-instance.json \
    --parameters @azure-parameters.json \
    --parameters containerImage="$ACR_LOGIN_SERVER/$IMAGE_NAME:$TAG" \
    --output json)

# Extract outputs
FQDN=$(echo "$DEPLOYMENT_OUTPUT" | jq -r '.properties.outputs.fqdn.value')
SSE_ENDPOINT=$(echo "$DEPLOYMENT_OUTPUT" | jq -r '.properties.outputs.sseEndpoint.value')
HEALTH_ENDPOINT=$(echo "$DEPLOYMENT_OUTPUT" | jq -r '.properties.outputs.healthEndpoint.value')

echo ""
echo "✅ Deployment completed successfully!"
echo ""
echo "📊 Deployment Information:"
echo "  Resource Group: $RESOURCE_GROUP_NAME"
echo "  Container Registry: $ACR_NAME"
echo "  Container Instance: $CONTAINER_NAME"
echo ""
echo "🌐 Endpoints:"
echo "  Public FQDN: $FQDN"
echo "  SSE Endpoint: $SSE_ENDPOINT"
echo "  Health Check: $HEALTH_ENDPOINT"
echo ""
echo "🎯 For Copilot Studio, use this SSE endpoint:"
echo "  $SSE_ENDPOINT"
echo ""
echo "🔍 To test the deployment:"
echo "  curl $HEALTH_ENDPOINT"
echo ""
echo "📝 To view logs:"
echo "  az container logs --resource-group $RESOURCE_GROUP_NAME --name $CONTAINER_NAME"
echo ""
echo "🗑️  To cleanup resources:"
echo "  az group delete --name $RESOURCE_GROUP_NAME --yes --no-wait"
