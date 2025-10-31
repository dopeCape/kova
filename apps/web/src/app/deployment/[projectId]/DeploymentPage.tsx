"use client";
import React, { useState, useEffect, useRef } from "react";
import { useSession } from "next-auth/react";
import {
  Loader2,
  CheckCircle,
  XCircle,
  Package,
  Rocket,
  Server,
  ExternalLink,
} from "lucide-react";

interface DeploymentStatusProps {
  projectId: string;
}

type DeploymentStatus =
  | "pending"
  | "building"
  | "deploying"
  | "deployed"
  | "failed";

const DeploymentStatus: React.FC<DeploymentStatusProps> = ({ projectId }) => {
  const { data: session, status: authStatus } = useSession();
  const [status, setStatus] = useState<DeploymentStatus>("pending");
  const [project, setProject] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const wsRef = useRef<WebSocket | null>(null);

  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const wsUrl = apiUrl.replace("http", "ws");

  // Fetch project details
  const fetchProject = async () => {
    if (!session?.user?.id) return;

    try {
      const response = await fetch(
        `${apiUrl}/api/v1/users/${session.user.id}/projects/${projectId}`,
        {
          headers: {
            Authorization: `Bearer ${session.accessToken}`,
            "Content-Type": "application/json",
          },
        },
      );

      if (response.ok) {
        const data = await response.json();
        setProject(data.project);
        setStatus(data.project.deployment_status);
      }
    } catch (error) {
      console.error("Error fetching project:", error);
    } finally {
      setLoading(false);
    }
  };

  // Connect to WebSocket
  useEffect(() => {
    if (authStatus !== "authenticated" || !session?.user?.id) return;

    fetchProject();

    // Connect WebSocket
    const ws = new WebSocket(
      `${wsUrl}/api/v1/users/${session.user.id}/projects/${projectId}/ws`,
    );
    wsRef.current = ws;

    ws.onopen = () => {
      console.log("WebSocket connected");
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "deployment_status" && data.status) {
          setStatus(data.status as DeploymentStatus);
        }
      } catch (error) {
        console.error("Error parsing WebSocket message:", error);
      }
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    ws.onclose = () => {
      console.log("WebSocket disconnected");
    };

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [authStatus, session?.user?.id, projectId]);

  if (authStatus === "loading" || loading) {
    return (
      <div className="w-full h-screen flex items-center justify-center bg-black">
        <div className="flex flex-col items-center space-y-4">
          <Loader2 className="w-12 h-12 animate-spin text-emerald-400" />
          <div className="text-zinc-300 font-medium">Loading deployment...</div>
        </div>
      </div>
    );
  }

  if (authStatus === "unauthenticated") {
    return (
      <div className="w-full h-screen flex items-center justify-center bg-black">
        <div className="text-center">
          <div className="text-zinc-400 font-mono text-sm mb-4">
            Please sign in to continue
          </div>
        </div>
      </div>
    );
  }

  const getStatusColor = () => {
    switch (status) {
      case "pending":
        return "text-yellow-400";
      case "building":
        return "text-blue-400";
      case "deploying":
        return "text-purple-400";
      case "deployed":
        return "text-emerald-400";
      case "failed":
        return "text-red-400";
      default:
        return "text-zinc-400";
    }
  };

  const getStatusIcon = () => {
    switch (status) {
      case "pending":
        return <Loader2 className="w-12 h-12 animate-spin" />;
      case "building":
        return <Package className="w-12 h-12 animate-pulse" />;
      case "deploying":
        return <Rocket className="w-12 h-12 animate-pulse" />;
      case "deployed":
        return <CheckCircle className="w-12 h-12" />;
      case "failed":
        return <XCircle className="w-12 h-12" />;
      default:
        return <Server className="w-12 h-12" />;
    }
  };

  const getStatusMessage = () => {
    switch (status) {
      case "pending":
        return "Queued for deployment";
      case "building":
        return "Building your application...";
      case "deploying":
        return "Deploying to production...";
      case "deployed":
        return "Successfully deployed!";
      case "failed":
        return "Deployment failed";
      default:
        return "Processing...";
    }
  };

  return (
    <div className="w-full min-h-screen bg-black p-6">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-white mb-2">
            {project?.name || "Project"}
          </h1>
          <p className="text-zinc-400 font-mono text-sm">
            {project?.repo_full_name}
          </p>
        </div>

        {/* Status Card */}
        <div className="bg-zinc-900/50 backdrop-blur-xl border border-zinc-800/50 rounded-xl p-8">
          <div className="flex flex-col items-center text-center space-y-6">
            {/* Status Icon */}
            <div className={`${getStatusColor()}`}>{getStatusIcon()}</div>

            {/* Status Message */}
            <div>
              <h2 className={`text-2xl font-bold mb-2 ${getStatusColor()}`}>
                {getStatusMessage()}
              </h2>
              <p className="text-zinc-400 font-mono text-sm">
                {status === "deployed"
                  ? "Your application is now live!"
                  : "Please wait while we process your deployment"}
              </p>
            </div>

            {/* Progress Steps */}
            <div className="w-full max-w-md mt-8">
              <div className="space-y-4">
                {/* Pending/Building */}
                <div className="flex items-center space-x-3">
                  <div
                    className={`w-8 h-8 rounded-full flex items-center justify-center ${status === "pending" ||
                        status === "building" ||
                        status === "deploying" ||
                        status === "deployed"
                        ? "bg-emerald-500"
                        : "bg-zinc-700"
                      }`}
                  >
                    {status === "building" ? (
                      <Loader2 className="w-4 h-4 animate-spin text-white" />
                    ) : (
                      <CheckCircle className="w-4 h-4 text-white" />
                    )}
                  </div>
                  <div className="flex-1">
                    <div className="text-sm font-medium text-white">
                      Building
                    </div>
                    <div className="text-xs text-zinc-500">
                      Compiling your application
                    </div>
                  </div>
                </div>

                {/* Deploying */}
                <div className="flex items-center space-x-3">
                  <div
                    className={`w-8 h-8 rounded-full flex items-center justify-center ${status === "deploying" || status === "deployed"
                        ? "bg-emerald-500"
                        : status === "building"
                          ? "bg-blue-500"
                          : "bg-zinc-700"
                      }`}
                  >
                    {status === "deploying" ? (
                      <Loader2 className="w-4 h-4 animate-spin text-white" />
                    ) : status === "deployed" ? (
                      <CheckCircle className="w-4 h-4 text-white" />
                    ) : (
                      <Rocket className="w-4 h-4 text-white" />
                    )}
                  </div>
                  <div className="flex-1">
                    <div className="text-sm font-medium text-white">
                      Deploying
                    </div>
                    <div className="text-xs text-zinc-500">
                      Starting your services
                    </div>
                  </div>
                </div>

                {/* Deployed */}
                <div className="flex items-center space-x-3">
                  <div
                    className={`w-8 h-8 rounded-full flex items-center justify-center ${status === "deployed" ? "bg-emerald-500" : "bg-zinc-700"
                      }`}
                  >
                    <CheckCircle className="w-4 h-4 text-white" />
                  </div>
                  <div className="flex-1">
                    <div className="text-sm font-medium text-white">Live</div>
                    <div className="text-xs text-zinc-500">
                      Application is running
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Success Actions */}
            {status === "deployed" && project?.domain && (
              <div className="mt-8 space-y-4 w-full max-w-md">
                <a
                  href={`https://${project.domain}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center justify-center space-x-2 w-full px-6 py-3 bg-emerald-500 hover:bg-emerald-400 text-white font-bold rounded-lg transition-all duration-200"
                >
                  <ExternalLink className="w-4 h-4" />
                  <span>Visit Your Application</span>
                </a>
                <button
                  onClick={() =>
                    (window.location.href = `/projects/${projectId}`)
                  }
                  className="flex items-center justify-center space-x-2 w-full px-6 py-3 bg-zinc-800 hover:bg-zinc-700 text-white rounded-lg transition-all duration-200"
                >
                  <span>View Project Details</span>
                </button>
              </div>
            )}

            {/* Failed State */}
            {status === "failed" && (
              <div className="mt-8 space-y-4 w-full max-w-md">
                <div className="text-center text-zinc-400 text-sm">
                  Something went wrong during deployment. Please try again or
                  contact support.
                </div>
                <button
                  onClick={() => (window.location.href = `/project/new`)}
                  className="flex items-center justify-center space-x-2 w-full px-6 py-3 bg-zinc-800 hover:bg-zinc-700 text-white rounded-lg transition-all duration-200"
                >
                  <span>Create New Project</span>
                </button>
              </div>
            )}
          </div>
        </div>

        {/* Project Info */}
        {project && (
          <div className="mt-6 bg-zinc-900/30 backdrop-blur-xl border border-zinc-800/50 rounded-xl p-6">
            <h3 className="text-sm font-medium text-zinc-300 mb-4 font-mono">
              Project Information
            </h3>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <div className="text-xs text-zinc-500 mb-1">Branch</div>
                <div className="text-sm text-white font-mono">
                  {project.repo_branch}
                </div>
              </div>
              {project.domain && (
                <div>
                  <div className="text-xs text-zinc-500 mb-1">Domain</div>
                  <div className="text-sm text-white font-mono">
                    {project.domain}
                  </div>
                </div>
              )}
              {project.port && (
                <div>
                  <div className="text-xs text-zinc-500 mb-1">Port</div>
                  <div className="text-sm text-white font-mono">
                    {project.port}
                  </div>
                </div>
              )}
              <div>
                <div className="text-xs text-zinc-500 mb-1">Status</div>
                <div className="text-sm text-white font-mono capitalize">
                  {status}
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default DeploymentStatus;
