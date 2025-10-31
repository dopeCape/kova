"use client";
import React, { useState, useEffect } from "react";
import {
  ChevronRight,
  Github,
  Search,
  Eye,
  Plus,
  Loader2,
  Star,
  Lock,
  GitBranch,
  Calendar,
  Package,
  Settings,
  Rocket,
  X,
} from "lucide-react";
import { useSession } from "next-auth/react";

interface GitHubAccount {
  id: string;
  github_username: string;
  github_id: number;
  avatar_url: string;
  created_at: string;
  updated_at: string;
}

interface Repository {
  id: string;
  name: string;
  full_name: string;
  private: boolean;
  updated_at: string;
  language: string;
  stars: number;
  description: string;
  default_branch: string;
  url: string;
}

interface RepositoryAnalysis {
  install: string[];
  build: string[];
  deploy: string;
  success: boolean;
}

interface EnvironmentVariable {
  id: string;
  key: string;
  value: string;
}

const GitHubAccountsWorkflow: React.FC = () => {
  const { data: session, status } = useSession();
  const [currentStep, setCurrentStep] = useState(1);
  const [accounts, setAccounts] = useState<GitHubAccount[]>([]);
  const [selectedAccount, setSelectedAccount] = useState<GitHubAccount | null>(
    null,
  );
  const [repositories, setRepositories] = useState<Repository[]>([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [accessToken, setAccessToken] = useState("");
  const [loading, setLoading] = useState(false);
  const [loadingRepos, setLoadingRepos] = useState(false);
  const [analyzingRepo, setAnalyzingRepo] = useState(false);
  const [error, setError] = useState("");
  const [repoError, setRepoError] = useState("");
  const [showTokenInput, setShowTokenInput] = useState(false);
  const [selectedRepository, setSelectedRepository] =
    useState<Repository | null>(null);
  const [analysis, setAnalysis] = useState<RepositoryAnalysis | null>(null);
  const [envVariables, setEnvVariables] = useState<EnvironmentVariable[]>([]);
  const [creatingProject, setCreatingProject] = useState(false);
  const [projectName, setProjectName] = useState("");
  const [domain, setDomain] = useState("");

  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  // Fetch user's GitHub accounts
  const fetchAccounts = async () => {
    if (!session?.user?.id) return;

    setLoading(true);
    try {
      const response = await fetch(
        `${apiUrl}/api/v1/users/${session.user.id}/accounts`,
        {
          headers: {
            Authorization: `Bearer ${session.accessToken}`,
            "Content-Type": "application/json",
          },
        },
      );

      if (response.ok) {
        const data = await response.json();
        setAccounts(data.accounts || []);
      } else {
        console.error("Failed to fetch accounts");
      }
    } catch (error) {
      console.error("Error fetching accounts:", error);
    } finally {
      setLoading(false);
    }
  };

  // Fetch repositories for a specific account
  const fetchRepositories = async (accountId: string) => {
    if (!session?.user?.id) return;

    setLoadingRepos(true);
    setRepoError("");
    try {
      const response = await fetch(
        `${apiUrl}/api/v1/users/${session.user.id}/${accountId}/repositories`,
        {
          headers: {
            Authorization: `Bearer ${session.accessToken}`,
            "Content-Type": "application/json",
          },
        },
      );

      if (response.ok) {
        const data = await response.json();
        setRepositories(data.repositories || []);
      } else {
        const errorData = await response.json();
        setRepoError(errorData.error || "Failed to fetch repositories");
      }
    } catch (error) {
      console.error("Error fetching repositories:", error);
      setRepoError(
        "Network error. Please check your connection and try again.",
      );
    } finally {
      setLoadingRepos(false);
    }
  };

  // Create new GitHub account
  const createAccount = async () => {
    if (!session?.user?.id || !accessToken) return;

    setLoading(true);
    setError("");

    try {
      const response = await fetch(
        `${apiUrl}/api/v1/users/${session.user.id}/accounts`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${session.accessToken}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            access_token: accessToken,
          }),
        },
      );

      if (response.ok) {
        const data = await response.json();
        setAccounts((prev) => [...prev, data.account]);
        setAccessToken("");
        setShowTokenInput(false);
        setError("");
      } else {
        const errorData = await response.json();
        setError(errorData.error || "Failed to add GitHub account");
      }
    } catch (error) {
      setError("Network error. Please try again.");
      console.error("Error creating account:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleAccountSelect = (account: GitHubAccount) => {
    setSelectedAccount(account);
    setCurrentStep(2);
    fetchRepositories(account.id);
  };

  // Analyze repository when selected
  const analyzeRepository = async (repo: Repository) => {
    if (!session?.user?.id || !selectedAccount) return;

    setSelectedRepository(repo);
    setAnalyzingRepo(true);
    setRepoError("");

    try {
      // Extract owner and repo name from fullName (format: owner/repo)
      const [owner, repoName] = repo.full_name.split("/");

      const response = await fetch(
        `${apiUrl}/api/v1/users/${session.user.id}/${selectedAccount.id}/repositorie/analyze`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${session.accessToken}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            repo_url: repo.url,
            branch: repo.default_branch,
            repo_id: parseInt(repo.id),
            repo_name: repoName,
            repo_owner: owner,
          }),
        },
      );

      if (response.ok) {
        const data = await response.json();
        setAnalysis(data.analysis);
        setCurrentStep(3);
      } else {
        const errorData = await response.json();
        setRepoError(errorData.error || "Failed to analyze repository");
      }
    } catch (error) {
      console.error("Error analyzing repository:", error);
      setRepoError("Network error. Please try again.");
    } finally {
      setAnalyzingRepo(false);
    }
  };

  const addEnvVariable = () => {
    setEnvVariables([
      ...envVariables,
      { id: Date.now().toString(), key: "", value: "" },
    ]);
  };

  const removeEnvVariable = (id: string) => {
    setEnvVariables(envVariables.filter((env) => env.id !== id));
  };

  const updateEnvVariable = (
    id: string,
    field: "key" | "value",
    value: string,
  ) => {
    setEnvVariables(
      envVariables.map((env) =>
        env.id === id ? { ...env, [field]: value } : env,
      ),
    );
  };

  // Create project with all data
  const createProject = async () => {
    if (!session?.user?.id || !selectedRepository || !selectedAccount) return;

    setCreatingProject(true);
    setError("");

    try {
      // Filter out empty env variables
      const validEnvVars = envVariables
        .filter((env) => env.key.trim() !== "" && env.value.trim() !== "")
        .map(({ key, value }) => ({ key, value }));

      // Extract owner and repo name from fullName
      const [owner, repoName] = selectedRepository.full_name.split("/");

      const response = await fetch(
        `${apiUrl}/api/v1/users/${session.user.id}/projects`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${session.accessToken}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            name: projectName || selectedRepository.name,
            domain: domain,
            repo_id: parseInt(selectedRepository.id),
            repo_name: repoName,
            repo_full_name: selectedRepository.full_name,
            repo_url: selectedRepository.url,
            repo_branch: selectedRepository.default_branch,
            env_variables: validEnvVars,
          }),
        },
      );

      if (response.ok) {
        const data = await response.json();
        // Redirect to deployment page
        window.location.href = `/deployment/${data.project.id}`;
      } else {
        const errorData = await response.json();
        setError(errorData.error || "Failed to create project");
      }
    } catch (error) {
      console.error("Error creating project:", error);
      setError("Network error. Please try again.");
    } finally {
      setCreatingProject(false);
    }
  };

  const getFilteredRepositories = () => {
    return repositories.filter(
      (repo) =>
        repo.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        repo.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
        repo.language.toLowerCase().includes(searchTerm.toLowerCase()),
    );
  };

  // Fetch accounts when component mounts and user is authenticated
  useEffect(() => {
    if (status === "authenticated" && session?.user?.id) {
      fetchAccounts();
    }
  }, [status, session?.user?.id]);

  if (status === "loading") {
    return (
      <div className="w-full h-full p-6">
        <div className="h-full bg-zinc-900/50 backdrop-blur-xl border border-zinc-800/50 rounded-xl p-6 flex flex-col">
          <div className="flex items-center justify-between mb-6">
            <div className="space-y-2">
              <div className="h-5 bg-zinc-800/50 rounded animate-pulse w-40"></div>
              <div className="h-3 bg-zinc-800/30 rounded animate-pulse w-56"></div>
            </div>
            <div className="h-9 bg-zinc-800/50 rounded animate-pulse w-28"></div>
          </div>

          <div className="flex-1 flex items-center justify-center">
            <div className="flex flex-col items-center space-y-4">
              <div className="relative">
                <div className="w-12 h-12 bg-zinc-800/50 rounded-full animate-pulse"></div>
                <Loader2 className="w-6 h-6 animate-spin text-emerald-400 absolute inset-0 m-auto" />
              </div>
              <div className="text-center space-y-2">
                <div className="text-zinc-300 font-medium">
                  Initializing workspace...
                </div>
                <div className="text-xs text-zinc-500 font-mono">
                  Loading your GitHub accounts
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (status === "unauthenticated") {
    return (
      <div className="w-full h-full flex items-center justify-center">
        <div className="text-center">
          <div className="text-zinc-400 font-mono text-sm mb-4">
            Please sign in to continue
          </div>
        </div>
      </div>
    );
  }

  // Step 1: GitHub Accounts List
  if (currentStep === 1) {
    return (
      <div className="w-full h-full p-6">
        <div className="h-full bg-zinc-900/50 backdrop-blur-xl border border-zinc-800/50 rounded-xl p-6 flex flex-col">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h2 className="text-lg font-bold text-white mb-1">
                GitHub Accounts
              </h2>
              <p className="text-xs text-zinc-400 font-mono">
                Select an account to deploy from
              </p>
            </div>
            <button
              onClick={() => setShowTokenInput(!showTokenInput)}
              className="flex items-center space-x-2 px-4 py-2 bg-emerald-500/10 hover:bg-emerald-500/20 border border-emerald-500/30 text-emerald-400 rounded-lg transition-all duration-200 text-sm"
            >
              <Plus className="w-4 h-4" />
              <span>Add Account</span>
            </button>
          </div>

          {/* Add Token Input */}
          {showTokenInput && (
            <div className="bg-zinc-900/30 border border-zinc-800/50 rounded-lg p-4 mb-4">
              <label className="block text-xs font-medium text-zinc-300 mb-2 font-mono">
                GitHub Personal Access Token
              </label>
              <div className="flex space-x-2">
                <input
                  type="password"
                  value={accessToken}
                  onChange={(e) => setAccessToken(e.target.value)}
                  placeholder="ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
                  className="flex-1 px-3 py-2 bg-zinc-900/80 border border-zinc-800/50 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500/50 transition-all duration-200 font-mono text-sm"
                />
                <button
                  onClick={createAccount}
                  disabled={!accessToken || loading}
                  className="px-4 py-2 bg-emerald-500 hover:bg-emerald-400 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-all duration-200 text-sm"
                >
                  {loading ? (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  ) : (
                    "Add"
                  )}
                </button>
              </div>
              {error && (
                <div className="mt-2 text-xs text-red-400 font-mono">
                  {error}
                </div>
              )}
              <div className="mt-2 text-xs text-zinc-500">
                Create a token at{" "}
                <a
                  href="https://github.com/settings/tokens"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-emerald-400 hover:underline"
                >
                  github.com/settings/tokens
                </a>
              </div>
            </div>
          )}

          {/* Loading State - Accounts Skeleton */}
          {loading && !showTokenInput && (
            <div className="flex-1 space-y-3 overflow-y-auto">
              {[...Array(3)].map((_, i) => (
                <div
                  key={i}
                  className="p-4 bg-zinc-900/30 border border-zinc-800/30 rounded-lg"
                >
                  <div className="flex items-center space-x-3">
                    <div className="w-10 h-10 bg-zinc-800/50 rounded-full animate-pulse"></div>
                    <div className="flex-1">
                      <div className="h-4 bg-zinc-800/50 rounded animate-pulse w-32 mb-2"></div>
                      <div className="h-3 bg-zinc-800/30 rounded animate-pulse w-48"></div>
                    </div>
                    <div className="w-5 h-5 bg-zinc-800/30 rounded animate-pulse"></div>
                  </div>
                </div>
              ))}
            </div>
          )}

          {/* No Accounts State */}
          {!loading && accounts.length === 0 && !showTokenInput && (
            <div className="flex-1 flex items-center justify-center">
              <div className="text-center">
                <Github className="w-12 h-12 text-zinc-600 mx-auto mb-4" />
                <div className="text-zinc-400 font-medium mb-2">
                  No GitHub accounts connected
                </div>
                <div className="text-xs text-zinc-500 mb-4">
                  Add your GitHub personal access token to get started
                </div>
                <button
                  onClick={() => setShowTokenInput(true)}
                  className="flex items-center space-x-2 px-4 py-2 bg-emerald-500 hover:bg-emerald-400 text-white rounded-lg transition-all duration-200 text-sm mx-auto"
                >
                  <Plus className="w-4 h-4" />
                  <span>Connect GitHub Account</span>
                </button>
              </div>
            </div>
          )}

          {/* Accounts List */}
          {!loading && accounts.length > 0 && (
            <div className="flex-1 space-y-3 overflow-y-auto">
              {accounts.map((account) => (
                <div
                  key={account.id}
                  onClick={() => handleAccountSelect(account)}
                  className="group p-4 bg-zinc-900/30 backdrop-blur-sm border border-zinc-800/30 hover:border-emerald-500/30 rounded-lg cursor-pointer transition-all duration-200 hover:bg-zinc-900/50"
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-3">
                      <img
                        src={account.avatar_url}
                        alt={account.github_username}
                        className="w-10 h-10 rounded-full border border-zinc-700/50"
                      />
                      <div>
                        <div className="flex items-center space-x-2">
                          <span className="text-white font-medium">
                            {account.github_username}
                          </span>
                          <Github className="w-4 h-4 text-zinc-400" />
                        </div>
                        <div className="text-xs text-zinc-400 font-mono">
                          Connected{" "}
                          {new Date(account.created_at).toLocaleDateString()}
                        </div>
                      </div>
                    </div>
                    <ChevronRight className="w-5 h-5 text-zinc-500 group-hover:text-emerald-400 transition-colors" />
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    );
  }

  // Step 2: Repository Selection
  if (currentStep === 2 && selectedAccount) {
    return (
      <div className="w-full h-full p-6">
        <div className="h-full bg-zinc-900/50 backdrop-blur-xl border border-zinc-800/50 rounded-xl p-6 flex flex-col">
          {/* Header */}
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center space-x-3">
              <button
                onClick={() => setCurrentStep(1)}
                className="p-2 hover:bg-zinc-800/50 rounded-lg transition-colors"
              >
                <ChevronRight className="w-4 h-4 text-zinc-400 transform rotate-180" />
              </button>
              <img
                src={selectedAccount.avatar_url}
                alt={selectedAccount.github_username}
                className="w-8 h-8 rounded-full border border-zinc-700/50"
              />
              <div>
                <h2 className="text-lg font-bold text-white">
                  {selectedAccount.github_username}
                </h2>
                <p className="text-xs text-zinc-400 font-mono">
                  Select a repository to deploy
                </p>
              </div>
            </div>
          </div>

          {/* Search */}
          <div className="relative mb-4">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-zinc-500" />
            <input
              type="text"
              placeholder="Search repositories..."
              className="w-full pl-10 pr-4 py-3 bg-zinc-900/80 backdrop-blur-sm border border-zinc-800/50 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500/50 transition-all duration-200 font-mono text-sm"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </div>

          {/* Repositories List */}
          <div className="flex-1 overflow-y-auto">
            {loadingRepos ? (
              <div className="space-y-3">
                {[...Array(5)].map((_, i) => (
                  <div
                    key={i}
                    className="p-4 bg-zinc-900/30 border border-zinc-800/30 rounded-lg"
                  >
                    <div className="flex items-center space-x-3">
                      <div className="w-2 h-2 bg-zinc-800/50 rounded-full animate-pulse"></div>
                      <div className="flex-1">
                        <div className="flex items-center space-x-2 mb-2">
                          <div className="h-4 bg-zinc-800/50 rounded animate-pulse w-32"></div>
                          <div className="w-16 h-5 bg-zinc-800/30 rounded animate-pulse"></div>
                        </div>
                        <div className="h-3 bg-zinc-800/30 rounded animate-pulse w-full mb-2"></div>
                        <div className="flex items-center space-x-4">
                          <div className="h-3 bg-zinc-800/30 rounded animate-pulse w-12"></div>
                          <div className="h-3 bg-zinc-800/30 rounded animate-pulse w-16"></div>
                        </div>
                      </div>
                      <div className="flex items-center space-x-2">
                        <div className="w-6 h-6 bg-zinc-800/30 rounded animate-pulse"></div>
                        <div className="w-4 h-4 bg-zinc-800/30 rounded animate-pulse"></div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : repoError ? (
              <div className="flex-1 flex items-center justify-center">
                <div className="text-center max-w-md">
                  <div className="w-16 h-16 bg-red-500/10 rounded-full flex items-center justify-center mx-auto mb-4">
                    <Github className="w-8 h-8 text-red-400" />
                  </div>
                  <div className="text-zinc-300 font-medium mb-2">
                    Failed to load repositories
                  </div>
                  <div className="text-xs text-zinc-400 mb-6 font-mono">
                    {repoError}
                  </div>
                  <button
                    onClick={() =>
                      selectedAccount && fetchRepositories(selectedAccount.id)
                    }
                    className="flex items-center space-x-2 px-4 py-2 bg-emerald-500 hover:bg-emerald-400 text-white rounded-lg transition-all duration-200 text-sm mx-auto"
                  >
                    <Loader2 className="w-4 h-4" />
                    <span>Try Again</span>
                  </button>
                </div>
              </div>
            ) : (
              <div className="space-y-2">
                {getFilteredRepositories().map((repo) => (
                  <div
                    key={repo.id}
                    onClick={() => analyzeRepository(repo)}
                    className="group p-4 bg-zinc-900/30 backdrop-blur-sm border border-zinc-800/30 hover:border-emerald-500/30 rounded-lg cursor-pointer transition-all duration-200 hover:bg-zinc-900/50"
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center space-x-3 flex-1">
                        <div className="w-2 h-2 bg-emerald-400 rounded-full"></div>
                        <div className="flex-1">
                          <div className="flex items-center space-x-2 mb-1">
                            <span className="text-white font-medium">
                              {repo.name}
                            </span>
                            {repo.private && (
                              <Lock className="w-3 h-3 text-zinc-500" />
                            )}
                            <span className="px-2 py-1 bg-zinc-700/50 text-zinc-300 text-xs rounded font-mono">
                              {repo.language}
                            </span>
                          </div>
                          <div className="text-xs text-zinc-400 mb-2">
                            {repo.description}
                          </div>
                          <div className="flex items-center space-x-4 text-xs text-zinc-500 font-mono">
                            <div className="flex items-center space-x-1">
                              <Star className="w-3 h-3" />
                              <span>{repo.stars}</span>
                            </div>
                            <div className="flex items-center space-x-1">
                              <GitBranch className="w-3 h-3" />
                              <span>{repo.default_branch}</span>
                            </div>
                            <div className="flex items-center space-x-1">
                              <Calendar className="w-3 h-3" />
                              <span>Updated {repo.updated_at}</span>
                            </div>
                          </div>
                        </div>
                      </div>
                      <div className="flex items-center space-x-2">
                        {analyzingRepo && selectedRepository?.id === repo.id ? (
                          <Loader2 className="w-4 h-4 animate-spin text-emerald-400" />
                        ) : (
                          <ChevronRight className="w-4 h-4 text-zinc-500 group-hover:text-emerald-400 transition-colors" />
                        )}
                      </div>
                    </div>
                  </div>
                ))}
                {!loadingRepos &&
                  !repoError &&
                  getFilteredRepositories().length === 0 && (
                    <div className="text-center py-12">
                      <Github className="w-12 h-12 text-zinc-600 mx-auto mb-4" />
                      <div className="text-zinc-400 font-medium mb-2">
                        No repositories found
                      </div>
                      <div className="text-xs text-zinc-500">
                        {searchTerm
                          ? "Try adjusting your search terms"
                          : "No repositories available for this account"}
                      </div>
                    </div>
                  )}
              </div>
            )}
          </div>
        </div>
      </div>
    );
  }

  // Step 3: Project Configuration
  if (currentStep === 3 && selectedRepository && analysis) {
    return (
      <div className="w-full h-full p-6">
        <div className="h-full bg-zinc-900/50 backdrop-blur-xl border border-zinc-800/50 rounded-xl p-6 flex flex-col">
          {/* Header */}
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center space-x-3">
              <button
                onClick={() => setCurrentStep(2)}
                className="p-2 hover:bg-zinc-800/50 rounded-lg transition-colors"
              >
                <ChevronRight className="w-4 h-4 text-zinc-400 transform rotate-180" />
              </button>
              <div>
                <h2 className="text-lg font-bold text-white">
                  {selectedRepository.name}
                </h2>
                <p className="text-xs text-zinc-400 font-mono">
                  Configure deployment settings
                </p>
              </div>
            </div>
          </div>

          <div className="flex-1 overflow-y-auto space-y-6">
            {/* Project Name */}
            <div className="bg-zinc-900/50 backdrop-blur-xl border border-zinc-800/50 rounded-xl p-4">
              <label className="block text-xs font-medium text-zinc-300 mb-2 font-mono">
                Project Name
              </label>
              <input
                type="text"
                value={projectName}
                onChange={(e) => setProjectName(e.target.value)}
                placeholder={selectedRepository.name}
                className="w-full px-3 py-2 bg-zinc-900/80 border border-zinc-800/50 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500/50 transition-all duration-200 font-mono text-sm"
              />
            </div>

            {/* Domain */}
            <div className="bg-zinc-900/50 backdrop-blur-xl border border-zinc-800/50 rounded-xl p-4">
              <label className="block text-xs font-medium text-zinc-300 mb-2 font-mono">
                Domain *
              </label>
              <input
                type="text"
                value={domain}
                onChange={(e) => setDomain(e.target.value)}
                placeholder="example.com"
                className="w-full px-3 py-2 bg-zinc-900/80 border border-zinc-800/50 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500/50 transition-all duration-200 font-mono text-sm"
              />
              <div className="mt-2 text-xs text-zinc-500">
                Enter the domain where your project will be accessible
              </div>
            </div>

            {/* Install Commands */}
            {analysis.install.length > 0 && (
              <div className="bg-zinc-900/50 backdrop-blur-xl border border-zinc-800/50 rounded-xl p-4">
                <div className="flex items-center space-x-2 mb-3">
                  <Package className="w-4 h-4 text-emerald-400" />
                  <h3 className="text-sm font-medium text-white">
                    Install Commands
                  </h3>
                </div>
                <div className="space-y-2">
                  {analysis.install.map((cmd, index) => (
                    <div
                      key={index}
                      className="bg-zinc-900/80 border border-zinc-800/50 rounded-lg p-3 font-mono text-sm text-zinc-300"
                    >
                      {cmd}
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Build Commands */}
            {analysis.build.length > 0 && (
              <div className="bg-zinc-900/50 backdrop-blur-xl border border-zinc-800/50 rounded-xl p-4">
                <div className="flex items-center space-x-2 mb-3">
                  <Settings className="w-4 h-4 text-blue-400" />
                  <h3 className="text-sm font-medium text-white">
                    Build Commands
                  </h3>
                </div>
                <div className="space-y-2">
                  {analysis.build.map((cmd, index) => (
                    <div
                      key={index}
                      className="bg-zinc-900/80 border border-zinc-800/50 rounded-lg p-3 font-mono text-sm text-zinc-300"
                    >
                      {cmd}
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Deploy Command */}
            {analysis.deploy && (
              <div className="bg-zinc-900/50 backdrop-blur-xl border border-zinc-800/50 rounded-xl p-4">
                <div className="flex items-center space-x-2 mb-3">
                  <Rocket className="w-4 h-4 text-purple-400" />
                  <h3 className="text-sm font-medium text-white">
                    Start Command
                  </h3>
                </div>
                <div className="bg-zinc-900/80 border border-zinc-800/50 rounded-lg p-3 font-mono text-sm text-zinc-300">
                  {analysis.deploy}
                </div>
              </div>
            )}

            {/* Environment Variables */}
            <div className="bg-zinc-900/50 backdrop-blur-xl border border-zinc-800/50 rounded-xl p-4">
              <div className="flex items-center space-x-2 mb-3">
                <Settings className="w-4 h-4 text-yellow-400" />
                <h3 className="text-sm font-medium text-white">
                  Environment Variables
                </h3>
              </div>

              {envVariables.length > 0 && (
                <div className="space-y-2 mb-4">
                  {envVariables.map((env) => (
                    <div
                      key={env.id}
                      className="flex items-center space-x-2 group"
                    >
                      <input
                        type="text"
                        value={env.key}
                        onChange={(e) =>
                          updateEnvVariable(env.id, "key", e.target.value)
                        }
                        placeholder="KEY"
                        className="flex-1 px-3 py-2 bg-zinc-900/80 border border-zinc-800/50 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500/50 transition-all duration-200 font-mono text-sm"
                      />
                      <input
                        type="text"
                        value={env.value}
                        onChange={(e) =>
                          updateEnvVariable(env.id, "value", e.target.value)
                        }
                        placeholder="value"
                        className="flex-1 px-3 py-2 bg-zinc-900/80 border border-zinc-800/50 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500/50 transition-all duration-200 font-mono text-sm"
                      />
                      <button
                        onClick={() => removeEnvVariable(env.id)}
                        className="p-2 text-zinc-500 hover:text-red-400 transition-colors duration-200"
                      >
                        <X className="w-4 h-4" />
                      </button>
                    </div>
                  ))}
                </div>
              )}

              <button
                onClick={addEnvVariable}
                className="flex items-center justify-center space-x-2 w-full p-3 bg-zinc-900/50 hover:bg-zinc-900/80 border border-zinc-800/50 hover:border-emerald-500/30 text-zinc-300 hover:text-white rounded-lg transition-all duration-200"
              >
                <Plus className="w-4 h-4" />
                <span className="font-mono text-xs">Add Variable</span>
              </button>
            </div>
          </div>

          {/* Actions */}
          <div className="flex flex-col space-y-4 mt-6 pt-4 border-t border-zinc-800/50">
            {error && (
              <div className="text-xs text-red-400 font-mono text-center">
                {error}
              </div>
            )}
            <div className="flex justify-between">
              <button
                onClick={() => setCurrentStep(2)}
                disabled={creatingProject}
                className="px-4 py-2 bg-zinc-800/50 hover:bg-zinc-700/50 text-white rounded-lg transition-all duration-200 text-sm disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Back
              </button>
              <button
                onClick={createProject}
                disabled={creatingProject || !domain.trim()}
                className="flex items-center space-x-2 px-6 py-2 bg-emerald-500 hover:bg-emerald-400 text-white font-bold rounded-lg transition-all duration-200 text-sm shadow-lg shadow-emerald-500/25 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {creatingProject && (
                  <Loader2 className="w-4 h-4 animate-spin" />
                )}
                <span>
                  {creatingProject ? "Creating..." : "Create Project"}
                </span>
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return null;
};

export default GitHubAccountsWorkflow;
