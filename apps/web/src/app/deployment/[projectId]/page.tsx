import DeploymentStatus from "./DeploymentPage";

interface PageProps {
  params: Promise<{
    projectId: string;
  }>;
}

export default async function DeploymentPage({ params }: PageProps) {
  const { projectId } = await params;

  return <DeploymentStatus projectId={projectId} />;
}
