"use client"
import React, { useState, useEffect } from 'react';
import {
  LayoutDashboard,
  Settings,
  FolderOpen,
  Rocket,
  Activity,
  Users,
  Database,
  GitBranch,
  Clock,
  CheckCircle,
  XCircle,
  AlertCircle,
  BarChart3,
  Code2,
  Server,
  Globe,
  Bell,
  User,
  Key,
  Shield,
  Trash2,
  Plus,
  Search,
  Filter,
  MoreVertical,
  ExternalLink,
  Eye,
  Play,
  Pause,
  RefreshCw,
  TrendingUp,
  TrendingDown,
  Wifi,
  WifiOff,
  Zap,
  Download,
  Upload,
  Monitor,
  HardDrive,
  Cpu,
  MemoryStick,
  Calendar,
  MapPin,
  FileText,
  Link,
  Copy,
  Check,
  AlertTriangle,
  Info,
  Github,
  GitlabIcon as Gitlab,
  Chrome,
  Smartphone,
  Tablet,
  Laptop
} from 'lucide-react';

const KovaDashboard = () => {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [currentTime, setCurrentTime] = useState(new Date());
  const [isOnline, setIsOnline] = useState(true);
  const [copied, setCopied] = useState('');

  useEffect(() => {
    const timer = setInterval(() => setCurrentTime(new Date()), 1000);
    return () => clearInterval(timer);
  }, []);

  const copyToClipboard = (text: string, key: string) => {
    navigator.clipboard.writeText(text);
    setCopied(key);
    setTimeout(() => setCopied(''), 2000);
  };

  const navigation = [
    { id: 'dashboard', name: 'Dashboard', icon: LayoutDashboard },
    { id: 'projects', name: 'Projects', icon: FolderOpen },
    { id: 'deployments', name: 'Deployments', icon: Rocket },
    { id: 'settings', name: 'Settings', icon: Settings },
  ];

  const Sidebar = () => (
    <div className="w-64 bg-zinc-900/50 backdrop-blur-xl border-r border-zinc-800/50 min-h-screen relative">
      <div className="absolute inset-0 bg-gradient-to-b from-emerald-950/10 to-transparent"></div>

      <div className="relative z-10 p-6">
        <div className="flex items-center space-x-3 mb-8">
          <div className="w-8 h-8 bg-emerald-500 rounded-lg flex items-center justify-center">
            <Code2 className="w-5 h-5 text-white" />
          </div>
          <div className="flex flex-col">
            <span className="text-white font-bold text-xl">Kova</span>
            <span className="text-zinc-500 text-xs font-mono">v2.1.4</span>
          </div>
        </div>

        <div className="mb-4 p-3 bg-zinc-800/30 rounded-lg border border-zinc-700/50">
          <div className="flex items-center space-x-2 mb-1">
            {isOnline ? <Wifi className="w-3 h-3 text-emerald-400" /> : <WifiOff className="w-3 h-3 text-red-400" />}
            <span className="text-xs text-zinc-400">Status</span>
          </div>
          <div className="text-sm text-white font-medium">
            {isOnline ? 'All Systems Operational' : 'Connection Issues'}
          </div>
          <div className="text-xs text-zinc-500 font-mono">
            {currentTime.toLocaleTimeString()}
          </div>
        </div>

        <nav className="space-y-2">
          {navigation.map((item) => (
            <button
              key={item.id}
              onClick={() => setActiveTab(item.id)}
              className={`w-full flex items-center space-x-3 px-4 py-3 rounded-xl transition-all relative ${activeTab === item.id
                ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
                : 'text-zinc-400 hover:text-white hover:bg-zinc-800/50'
                }`}
            >
              <item.icon className="w-5 h-5" />
              <span className="font-medium">{item.name}</span>
              {activeTab === item.id && (
                <div className="absolute right-2 w-1 h-6 bg-emerald-500 rounded-full"></div>
              )}
            </button>
          ))}
        </nav>

        <div className="mt-8 p-4 bg-zinc-900/50 border border-zinc-800/50 rounded-xl">
          <div className="flex items-center space-x-2 mb-2">
            <Server className="w-4 h-4 text-emerald-400" />
            <span className="text-sm font-medium text-white">Infrastructure</span>
          </div>
          <div className="space-y-2 text-xs">
            <div className="flex justify-between">
              <span className="text-zinc-400">Regions</span>
              <span className="text-white">3 active</span>
            </div>
            <div className="flex justify-between">
              <span className="text-zinc-400">Edge Nodes</span>
              <span className="text-white">47 online</span>
            </div>
            <div className="flex justify-between">
              <span className="text-zinc-400">Uptime</span>
              <span className="text-emerald-400">99.98%</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );

  const DashboardPage = () => (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white mb-2">Dashboard</h1>
          <div className="flex items-center space-x-4 text-zinc-400">
            <span>Welcome back, Sarah Chen</span>
            <div className="flex items-center space-x-1">
              <MapPin className="w-3 h-3" />
              <span className="text-xs">San Francisco, CA</span>
            </div>
            <div className="flex items-center space-x-1">
              <Calendar className="w-3 h-3" />
              <span className="text-xs">{currentTime.toLocaleDateString()}</span>
            </div>
          </div>
        </div>
        <div className="flex items-center space-x-3">
          <div className="bg-zinc-900/50 border border-zinc-800/50 rounded-lg px-3 py-2">
            <div className="flex items-center space-x-2">
              <div className="w-2 h-2 bg-emerald-500 rounded-full animate-pulse"></div>
              <span className="text-sm text-white">Live</span>
            </div>
          </div>
          <button className="bg-emerald-500 hover:bg-emerald-400 text-white px-6 py-3 rounded-xl font-medium transition-all transform hover:scale-105 flex items-center space-x-2 shadow-lg shadow-emerald-500/25">
            <Plus className="w-4 h-4" />
            <span>New Project</span>
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {[
          {
            label: 'Active Projects',
            value: '27',
            change: '+3 this week',
            trend: 'up',
            icon: FolderOpen,
            color: 'emerald',
            details: '4 deployed today'
          },
          {
            label: 'Total Deployments',
            value: '1,247',
            change: '+89 this month',
            trend: 'up',
            icon: Rocket,
            color: 'blue',
            details: 'Avg 2.3/day'
          },
          {
            label: 'Success Rate',
            value: '99.7%',
            change: '+0.2% vs last week',
            trend: 'up',
            icon: CheckCircle,
            color: 'green',
            details: '3 failures this month'
          },
          {
            label: 'Avg Deploy Time',
            value: '1m 23s',
            change: '-12s improvement',
            trend: 'up',
            icon: Clock,
            color: 'purple',
            details: 'Fastest: 28s'
          },
        ].map((stat, index) => (
          <div key={index} className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6 relative overflow-hidden group hover:border-emerald-500/30 transition-all">
            <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <div className="relative z-10">
              <div className="flex items-center justify-between mb-3">
                <stat.icon className="w-8 h-8 text-emerald-400" />
                <div className={`flex items-center space-x-1 text-xs ${stat.trend === 'up' ? 'text-emerald-400' : 'text-red-400'}`}>
                  {stat.trend === 'up' ? <TrendingUp className="w-3 h-3" /> : <TrendingDown className="w-3 h-3" />}
                  <span>{stat.change}</span>
                </div>
              </div>
              <p className="text-zinc-400 text-sm font-medium mb-1">{stat.label}</p>
              <p className="text-3xl font-bold text-white mb-1">{stat.value}</p>
              <p className="text-zinc-500 text-xs">{stat.details}</p>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-xl font-semibold text-white flex items-center space-x-2">
              <Activity className="w-5 h-5 text-emerald-400" />
              <span>Recent Activity</span>
            </h3>
            <button className="text-zinc-400 hover:text-emerald-400 transition-colors text-sm">View All</button>
          </div>
          <div className="space-y-4">
            {[
              {
                action: 'Deployed successfully',
                project: 'ecommerce-frontend',
                time: '2 minutes ago',
                status: 'success',
                user: 'Sarah Chen',
                commit: 'fix: update checkout flow validation',
                commitHash: 'a7f3c92',
                duration: '1m 45s',
                url: 'ecommerce-frontend-git-main-acme.vercel.app'
              },
              {
                action: 'Build failed',
                project: 'api-gateway',
                time: '15 minutes ago',
                status: 'error',
                user: 'Mike Johnson',
                commit: 'feat: add rate limiting middleware',
                commitHash: 'b8e4d31',
                duration: '3m 12s',
                error: 'TypeScript compilation failed'
              },
              {
                action: 'Deployed successfully',
                project: 'marketing-site',
                time: '1 hour ago',
                status: 'success',
                user: 'Emily Rodriguez',
                commit: 'content: update pricing page',
                commitHash: 'c9f5e82',
                duration: '52s',
                url: 'marketing-git-main-acme.vercel.app'
              },
              {
                action: 'Building',
                project: 'admin-dashboard',
                time: '2 hours ago',
                status: 'pending',
                user: 'David Kim',
                commit: 'refactor: migrate to React 18',
                commitHash: 'd1a6f93',
                duration: '4m 33s'
              },
              {
                action: 'Deployed successfully',
                project: 'blog-cms',
                time: '3 hours ago',
                status: 'success',
                user: 'Lisa Wang',
                commit: 'feat: add image optimization',
                commitHash: 'e2b7g04',
                duration: '1m 18s',
                url: 'blog-cms-git-main-acme.vercel.app'
              },
            ].map((activity, index) => (
              <div key={index} className="flex items-start space-x-4 p-4 rounded-lg hover:bg-zinc-800/30 transition-colors group">
                <div className={`w-2 h-2 rounded-full mt-2 ${activity.status === 'success' ? 'bg-emerald-500' :
                  activity.status === 'error' ? 'bg-red-500' : 'bg-yellow-500'
                  }`}></div>
                <div className="flex-1 min-w-0">
                  <div className="flex items-center space-x-2 mb-1">
                    <p className="text-white font-medium">
                      {activity.action} <span className="text-emerald-400 font-mono">{activity.project}</span>
                    </p>
                    {activity.status === 'success' && activity.url && (
                      <button className="opacity-0 group-hover:opacity-100 transition-opacity">
                        <ExternalLink className="w-3 h-3 text-zinc-400 hover:text-emerald-400" />
                      </button>
                    )}
                  </div>
                  <div className="flex items-center space-x-4 text-xs text-zinc-400">
                    <span>by {activity.user}</span>
                    <span>•</span>
                    <span className="font-mono">{activity.commitHash}</span>
                    <span>•</span>
                    <span>{activity.duration}</span>
                    <span>•</span>
                    <span>{activity.time}</span>
                  </div>
                  <p className="text-xs text-zinc-500 mt-1 truncate">{activity.commit}</p>
                  {activity.error && (
                    <div className="flex items-center space-x-1 mt-2">
                      <AlertTriangle className="w-3 h-3 text-red-400" />
                      <span className="text-xs text-red-400">{activity.error}</span>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="space-y-6">
          <div className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6">
            <h3 className="text-xl font-semibold text-white mb-6 flex items-center space-x-2">
              <BarChart3 className="w-5 h-5 text-emerald-400" />
              <span>System Resources</span>
            </h3>
            <div className="space-y-4">
              {[
                { label: 'CPU Usage', value: 47, unit: '%', icon: Cpu, color: 'emerald', details: '2.4 GHz avg' },
                { label: 'Memory', value: 62, unit: '%', icon: MemoryStick, color: 'blue', details: '12.4GB / 20GB' },
                { label: 'Storage', value: 34, unit: '%', icon: HardDrive, color: 'purple', details: '340GB / 1TB' },
                { label: 'Network I/O', value: 18, unit: 'Mbps', icon: Monitor, color: 'orange', details: '↑ 12 ↓ 45' },
              ].map((resource, index) => (
                <div key={index}>
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center space-x-2">
                      <resource.icon className="w-4 h-4 text-zinc-400" />
                      <span className="text-zinc-400 text-sm">{resource.label}</span>
                    </div>
                    <div className="text-right">
                      <span className="text-white font-mono text-sm">{resource.value}{resource.unit}</span>
                      <p className="text-xs text-zinc-500">{resource.details}</p>
                    </div>
                  </div>
                  <div className="w-full bg-zinc-800 rounded-full h-2">
                    <div
                      className={`h-2 rounded-full ${resource.color === 'emerald' ? 'bg-emerald-500' :
                        resource.color === 'blue' ? 'bg-blue-500' :
                          resource.color === 'purple' ? 'bg-purple-500' : 'bg-orange-500'
                        }`}
                      style={{ width: `${resource.value}%` }}
                    ></div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6">
            <h3 className="text-xl font-semibold text-white mb-4 flex items-center space-x-2">
              <Globe className="w-5 h-5 text-emerald-400" />
              <span>Global Traffic</span>
            </h3>
            <div className="space-y-3">
              {[
                { region: 'North America', requests: '2.4M', latency: '45ms', color: 'emerald' },
                { region: 'Europe', requests: '1.8M', latency: '62ms', color: 'blue' },
                { region: 'Asia Pacific', requests: '3.1M', latency: '78ms', color: 'purple' },
                { region: 'South America', requests: '456K', latency: '92ms', color: 'orange' },
              ].map((region, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-zinc-800/30 rounded-lg">
                  <div className="flex items-center space-x-3">
                    <div className={`w-3 h-3 rounded-full ${region.color === 'emerald' ? 'bg-emerald-500' :
                      region.color === 'blue' ? 'bg-blue-500' :
                        region.color === 'purple' ? 'bg-purple-500' : 'bg-orange-500'
                      }`}></div>
                    <span className="text-white text-sm">{region.region}</span>
                  </div>
                  <div className="text-right">
                    <div className="text-white font-mono text-sm">{region.requests}</div>
                    <div className="text-zinc-400 text-xs">{region.latency}</div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );

  const ProjectsPage = () => (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white mb-2">Projects</h1>
          <p className="text-zinc-400">Manage your deployment projects across environments</p>
        </div>
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2 bg-zinc-900/50 border border-zinc-800/50 rounded-lg px-3 py-2">
            <Filter className="w-4 h-4 text-zinc-400" />
            <select className="bg-transparent text-white text-sm focus:outline-none">
              <option value="all">All Projects</option>
              <option value="active">Active Only</option>
              <option value="errors">With Errors</option>
              <option value="building">Building</option>
            </select>
          </div>
          <div className="relative">
            <Search className="w-4 h-4 absolute left-3 top-1/2 transform -translate-y-1/2 text-zinc-400" />
            <input
              type="text"
              placeholder="Search projects..."
              className="pl-10 pr-4 py-2 bg-zinc-900/50 border border-zinc-800 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500 transition-colors w-64"
            />
          </div>
          <button className="bg-emerald-500 hover:bg-emerald-400 text-white px-6 py-2 rounded-lg font-medium transition-all flex items-center space-x-2 shadow-lg shadow-emerald-500/25">
            <Plus className="w-4 h-4" />
            <span>New Project</span>
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
        {[
          {
            name: 'ecommerce-frontend',
            status: 'active',
            deployments: 47,
            lastDeploy: '2 minutes ago',
            framework: 'Next.js 14',
            repository: 'github.com/acme/ecommerce-frontend',
            branch: 'main',
            domain: 'shop.acme.com',
            team: 'Frontend Team',
            visitors: '24.5K/day',
            environments: ['production', 'staging', 'preview']
          },
          {
            name: 'api-gateway',
            status: 'error',
            deployments: 156,
            lastDeploy: '15 minutes ago',
            framework: 'Node.js',
            repository: 'github.com/acme/api-gateway',
            branch: 'develop',
            domain: 'api.acme.com',
            team: 'Backend Team',
            requests: '1.2M/day',
            environments: ['production', 'staging']
          },
          {
            name: 'marketing-site',
            status: 'active',
            deployments: 23,
            lastDeploy: '1 hour ago',
            framework: 'Astro',
            repository: 'github.com/acme/marketing-site',
            branch: 'main',
            domain: 'acme.com',
            team: 'Marketing',
            visitors: '8.7K/day',
            environments: ['production', 'staging']
          },
          {
            name: 'admin-dashboard',
            status: 'building',
            deployments: 89,
            lastDeploy: '2 hours ago',
            framework: 'React 18',
            repository: 'github.com/acme/admin-dashboard',
            branch: 'feature/react-18-migration',
            domain: 'admin.acme.com',
            team: 'Platform Team',
            users: '45 active',
            environments: ['production', 'staging', 'development']
          },
          {
            name: 'blog-cms',
            status: 'active',
            deployments: 34,
            lastDeploy: '3 hours ago',
            framework: 'Nuxt.js',
            repository: 'github.com/acme/blog-cms',
            branch: 'main',
            domain: 'blog.acme.com',
            team: 'Content Team',
            articles: '245 published',
            environments: ['production', 'staging']
          },
          {
            name: 'mobile-app-api',
            status: 'inactive',
            deployments: 67,
            lastDeploy: '3 days ago',
            framework: 'Express.js',
            repository: 'github.com/acme/mobile-api',
            branch: 'main',
            domain: 'mobile-api.acme.com',
            team: 'Mobile Team',
            endpoints: '47 routes',
            environments: ['production', 'staging']
          },
        ].map((project, index) => (
          <div key={index} className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6 hover:border-emerald-500/30 transition-all group relative overflow-hidden">
            <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>

            <div className="relative z-10">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center space-x-3">
                  <div className={`w-3 h-3 rounded-full ${project.status === 'active' ? 'bg-emerald-500' :
                    project.status === 'error' ? 'bg-red-500' :
                      project.status === 'building' ? 'bg-yellow-500 animate-pulse' : 'bg-zinc-500'
                    }`}></div>
                  <div>
                    <h3 className="text-white font-semibold font-mono">{project.name}</h3>
                    <p className="text-xs text-zinc-400">{project.team}</p>
                  </div>
                </div>
                <div className="flex items-center space-x-2">
                  <div className="flex items-center space-x-1">
                    {project.environments.map((env, idx) => (
                      <div key={idx} className={`w-2 h-2 rounded-full ${env === 'production' ? 'bg-emerald-500' :
                        env === 'staging' ? 'bg-yellow-500' : 'bg-blue-500'
                        }`} title={env}></div>
                    ))}
                  </div>
                  <button className="text-zinc-400 hover:text-white transition-colors opacity-0 group-hover:opacity-100">
                    <MoreVertical className="w-4 h-4" />
                  </button>
                </div>
              </div>

              <div className="space-y-3 mb-4">
                <div className="flex items-center justify-between text-sm">
                  <span className="text-zinc-400">Framework</span>
                  <span className="text-white">{project.framework}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-zinc-400">Domain</span>
                  <div className="flex items-center space-x-1">
                    <span className="text-emerald-400 font-mono text-xs">{project.domain}</span>
                    <ExternalLink className="w-3 h-3 text-zinc-400 hover:text-emerald-400 cursor-pointer" />
                  </div>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-zinc-400">Deployments</span>
                  <span className="text-white">{project.deployments}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-zinc-400">Last Deploy</span>
                  <span className="text-white">{project.lastDeploy}</span>
                </div>
                <div className="flex items-center justify-between text-sm">
                  <span className="text-zinc-400">
                    {project.visitors ? 'Visitors' : project.requests ? 'Requests' : project.users ? 'Users' : project.articles ? 'Articles' : 'Endpoints'}
                  </span>
                  <span className="text-white">
                    {project.visitors || project.requests || project.users || project.articles || project.endpoints}
                  </span>
                </div>
              </div>

              <div className="flex items-center space-x-2 mb-4">
                <GitBranch className="w-3 h-3 text-zinc-400" />
                <span className="text-xs text-zinc-400 font-mono">{project.branch}</span>
                <Github className="w-3 h-3 text-zinc-400" />
              </div>

              <div className="flex items-center space-x-2">
                <button className="flex-1 bg-zinc-800 hover:bg-zinc-700 text-white py-2 px-4 rounded-lg transition-colors flex items-center justify-center space-x-2 text-sm">
                  <Eye className="w-4 h-4" />
                  <span>View</span>
                </button>
                <button className={`${project.status === 'building'
                  ? 'bg-zinc-700 cursor-not-allowed'
                  : 'bg-emerald-500 hover:bg-emerald-400'
                  } text-white py-2 px-4 rounded-lg transition-colors flex items-center space-x-2 text-sm`}>
                  {project.status === 'building' ? (
                    <>
                      <RefreshCw className="w-4 h-4 animate-spin" />
                      <span>Building</span>
                    </>
                  ) : (
                    <>
                      <Rocket className="w-4 h-4" />
                      <span>Deploy</span>
                    </>
                  )}
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold text-white">Quick Stats</h3>
          <span className="text-xs text-zinc-400">Last 30 days</span>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="text-center">
            <div className="text-2xl font-bold text-emerald-400">892</div>
            <div className="text-xs text-zinc-400">Total Deployments</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-blue-400">47.2M</div>
            <div className="text-xs text-zinc-400">Total Requests</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-purple-400">1.2TB</div>
            <div className="text-xs text-zinc-400">Data Transfer</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold text-orange-400">99.96%</div>
            <div className="text-xs text-zinc-400">Uptime</div>
          </div>
        </div>
      </div>
    </div>
  );

  const DeploymentsPage = () => (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white mb-2">Deployments</h1>
          <p className="text-zinc-400">Track deployment history, logs, and performance metrics</p>
        </div>
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2">
            <button className="bg-zinc-800 hover:bg-zinc-700 text-white px-4 py-2 rounded-lg transition-colors flex items-center space-x-2 text-sm">
              <Filter className="w-4 h-4" />
              <span>Filter</span>
            </button>
            <select className="bg-zinc-900/50 border border-zinc-800 rounded-lg px-3 py-2 text-white text-sm focus:outline-none focus:border-emerald-500">
              <option value="all">All Status</option>
              <option value="success">Success</option>
              <option value="error">Error</option>
              <option value="building">Building</option>
            </select>
            <select className="bg-zinc-900/50 border border-zinc-800 rounded-lg px-3 py-2 text-white text-sm focus:outline-none focus:border-emerald-500">
              <option value="7">Last 7 days</option>
              <option value="30">Last 30 days</option>
              <option value="90">Last 90 days</option>
            </select>
          </div>
          <button className="bg-zinc-800 hover:bg-zinc-700 text-white px-4 py-2 rounded-lg transition-colors flex items-center space-x-2 text-sm">
            <RefreshCw className="w-4 h-4" />
            <span>Refresh</span>
          </button>
        </div>
      </div>

      <div className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-zinc-800/50">
              <tr>
                <th className="text-left py-4 px-6 text-zinc-300 font-medium">Status</th>
                <th className="text-left py-4 px-6 text-zinc-300 font-medium">Project</th>
                <th className="text-left py-4 px-6 text-zinc-300 font-medium">Branch</th>
                <th className="text-left py-4 px-6 text-zinc-300 font-medium">Commit</th>
                <th className="text-left py-4 px-6 text-zinc-300 font-medium">Duration</th>
                <th className="text-left py-4 px-6 text-zinc-300 font-medium">Deployed</th>
                <th className="text-left py-4 px-6 text-zinc-300 font-medium">Environment</th>
                <th className="text-left py-4 px-6 text-zinc-300 font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {[
                {
                  status: 'success',
                  project: 'ecommerce-frontend',
                  branch: 'main',
                  commit: 'fix: update checkout flow validation',
                  commitHash: 'a7f3c92',
                  duration: '1m 45s',
                  deployed: '2 minutes ago',
                  environment: 'production',
                  user: 'Sarah Chen',
                  url: 'ecommerce-frontend-git-main-acme.vercel.app',
                  size: '2.4 MB'
                },
                {
                  status: 'error',
                  project: 'api-gateway',
                  branch: 'develop',
                  commit: 'feat: add rate limiting middleware',
                  commitHash: 'b8e4d31',
                  duration: '3m 12s',
                  deployed: '15 minutes ago',
                  environment: 'staging',
                  user: 'Mike Johnson',
                  error: 'TypeScript compilation failed at src/middleware/rateLimit.ts:23',
                  size: '1.8 MB'
                },
                {
                  status: 'success',
                  project: 'marketing-site',
                  branch: 'main',
                  commit: 'content: update pricing page',
                  commitHash: 'c9f5e82',
                  duration: '52s',
                  deployed: '1 hour ago',
                  environment: 'production',
                  user: 'Emily Rodriguez',
                  url: 'marketing-git-main-acme.vercel.app',
                  size: '890 KB'
                },
                {
                  status: 'building',
                  project: 'admin-dashboard',
                  branch: 'feature/react-18-migration',
                  commit: 'refactor: migrate to React 18',
                  commitHash: 'd1a6f93',
                  duration: '4m 33s',
                  deployed: '2 hours ago',
                  environment: 'preview',
                  user: 'David Kim',
                  progress: '78%'
                },
                {
                  status: 'success',
                  project: 'blog-cms',
                  branch: 'main',
                  commit: 'feat: add image optimization',
                  commitHash: 'e2b7g04',
                  duration: '1m 18s',
                  deployed: '3 hours ago',
                  environment: 'production',
                  user: 'Lisa Wang',
                  url: 'blog-cms-git-main-acme.vercel.app',
                  size: '1.2 MB'
                },
                {
                  status: 'success',
                  project: 'mobile-app-api',
                  branch: 'hotfix/auth-security',
                  commit: 'security: fix JWT token validation',
                  commitHash: 'f3c8h15',
                  duration: '2m 34s',
                  deployed: '6 hours ago',
                  environment: 'production',
                  user: 'Alex Thompson',
                  url: 'mobile-api-git-hotfix-acme.vercel.app',
                  size: '3.1 MB'
                },
                {
                  status: 'cancelled',
                  project: 'admin-dashboard',
                  branch: 'feature/new-ui',
                  commit: 'ui: implement new design system',
                  commitHash: 'g4d9i26',
                  duration: '1m 02s',
                  deployed: '8 hours ago',
                  environment: 'preview',
                  user: 'Sarah Chen',
                  reason: 'Cancelled by user'
                },
              ].map((deployment, index) => (
                <tr key={index} className="border-t border-zinc-800/50 hover:bg-zinc-800/20 transition-colors group">
                  <td className="py-4 px-6">
                    <div className="flex items-center space-x-2">
                      {deployment.status === 'success' && <CheckCircle className="w-4 h-4 text-emerald-500" />}
                      {deployment.status === 'error' && <XCircle className="w-4 h-4 text-red-500" />}
                      {deployment.status === 'building' && <RefreshCw className="w-4 h-4 text-yellow-500 animate-spin" />}
                      {deployment.status === 'cancelled' && <AlertCircle className="w-4 h-4 text-orange-500" />}
                      <div>
                        <span className="text-white capitalize text-sm">{deployment.status}</span>
                        {deployment.progress && (
                          <div className="text-xs text-zinc-400">{deployment.progress}</div>
                        )}
                      </div>
                    </div>
                  </td>
                  <td className="py-4 px-6">
                    <div>
                      <span className="text-emerald-400 font-mono text-sm">{deployment.project}</span>
                      <div className="text-xs text-zinc-400">by {deployment.user}</div>
                    </div>
                  </td>
                  <td className="py-4 px-6">
                    <div className="flex items-center space-x-2">
                      <GitBranch className="w-3 h-3 text-zinc-400" />
                      <span className="text-white font-mono text-sm">{deployment.branch}</span>
                    </div>
                  </td>
                  <td className="py-4 px-6">
                    <div>
                      <div className="flex items-center space-x-2">
                        <span className="text-zinc-300 font-mono text-sm">{deployment.commitHash}</span>
                        <button
                          onClick={() => copyToClipboard(deployment.commitHash, `commit-${index}`)}
                          className="opacity-0 group-hover:opacity-100 transition-opacity"
                        >
                          {copied === `commit-${index}` ?
                            <Check className="w-3 h-3 text-emerald-400" /> :
                            <Copy className="w-3 h-3 text-zinc-400 hover:text-white" />
                          }
                        </button>
                      </div>
                      <div className="text-xs text-zinc-400 truncate max-w-xs" title={deployment.commit}>
                        {deployment.commit}
                      </div>
                    </div>
                  </td>
                  <td className="py-4 px-6">
                    <div>
                      <span className="text-white text-sm">{deployment.duration}</span>
                      {deployment.size && (
                        <div className="text-xs text-zinc-400">{deployment.size}</div>
                      )}
                    </div>
                  </td>
                  <td className="py-4 px-6">
                    <span className="text-zinc-400 text-sm">{deployment.deployed}</span>
                  </td>
                  <td className="py-4 px-6">
                    <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${deployment.environment === 'production' ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30' :
                      deployment.environment === 'staging' ? 'bg-yellow-500/20 text-yellow-400 border border-yellow-500/30' :
                        'bg-blue-500/20 text-blue-400 border border-blue-500/30'
                      }`}>
                      {deployment.environment}
                    </span>
                  </td>
                  <td className="py-4 px-6">
                    <div className="flex items-center space-x-2">
                      <button className="text-zinc-400 hover:text-emerald-400 transition-colors" title="View Logs">
                        <Eye className="w-4 h-4" />
                      </button>
                      {deployment.url && (
                        <button className="text-zinc-400 hover:text-emerald-400 transition-colors" title="Open Deployment">
                          <ExternalLink className="w-4 h-4" />
                        </button>
                      )}
                      <button className="text-zinc-400 hover:text-emerald-400 transition-colors" title="View Details">
                        <Info className="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6">
          <h3 className="text-lg font-semibold text-white mb-4 flex items-center space-x-2">
            <BarChart3 className="w-5 h-5 text-emerald-400" />
            <span>Deployment Frequency</span>
          </h3>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-zinc-400">Today</span>
              <div className="flex items-center space-x-2">
                <div className="w-32 bg-zinc-800 rounded-full h-2">
                  <div className="bg-emerald-500 h-2 rounded-full" style={{ width: '75%' }}></div>
                </div>
                <span className="text-white text-sm font-mono">12/16</span>
              </div>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-zinc-400">This week</span>
              <div className="flex items-center space-x-2">
                <div className="w-32 bg-zinc-800 rounded-full h-2">
                  <div className="bg-blue-500 h-2 rounded-full" style={{ width: '89%' }}></div>
                </div>
                <span className="text-white text-sm font-mono">89/100</span>
              </div>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-zinc-400">This month</span>
              <div className="flex items-center space-x-2">
                <div className="w-32 bg-zinc-800 rounded-full h-2">
                  <div className="bg-purple-500 h-2 rounded-full" style={{ width: '67%' }}></div>
                </div>
                <span className="text-white text-sm font-mono">267/400</span>
              </div>
            </div>
          </div>
        </div>

        <div className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6">
          <h3 className="text-lg font-semibold text-white mb-4 flex items-center space-x-2">
            <Clock className="w-5 h-5 text-emerald-400" />
            <span>Performance Metrics</span>
          </h3>
          <div className="grid grid-cols-2 gap-4">
            <div className="text-center p-3 bg-zinc-800/30 rounded-lg">
              <div className="text-lg font-bold text-emerald-400">1m 23s</div>
              <div className="text-xs text-zinc-400">Avg Build Time</div>
            </div>
            <div className="text-center p-3 bg-zinc-800/30 rounded-lg">
              <div className="text-lg font-bold text-blue-400">99.7%</div>
              <div className="text-xs text-zinc-400">Success Rate</div>
            </div>
            <div className="text-center p-3 bg-zinc-800/30 rounded-lg">
              <div className="text-lg font-bold text-purple-400">28s</div>
              <div className="text-xs text-zinc-400">Fastest Deploy</div>
            </div>
            <div className="text-center p-3 bg-zinc-800/30 rounded-lg">
              <div className="text-lg font-bold text-orange-400">2.1 MB</div>
              <div className="text-xs text-zinc-400">Avg Bundle Size</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );

  const SettingsPage = () => (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-white mb-2">Settings</h1>
        <p className="text-zinc-400">Manage your account, team, and deployment preferences</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <div className="lg:col-span-1">
          <div className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6 sticky top-6">
            <nav className="space-y-2">
              {[
                { id: 'profile', name: 'Profile', icon: User, badge: null },
                { id: 'security', name: 'Security', icon: Shield, badge: '2FA' },
                { id: 'api', name: 'API Keys', icon: Key, badge: '3' },
                { id: 'notifications', name: 'Notifications', icon: Bell, badge: 'New' },
                { id: 'integrations', name: 'Integrations', icon: Globe, badge: '5' },
                { id: 'billing', name: 'Billing', icon: FileText, badge: null },
                { id: 'team', name: 'Team', icon: Users, badge: '12' },
              ].map((item) => (
                <button
                  key={item.id}
                  className="w-full flex items-center justify-between px-4 py-3 rounded-lg text-zinc-400 hover:text-white hover:bg-zinc-800/50 transition-all group"
                >
                  <div className="flex items-center space-x-3">
                    <item.icon className="w-4 h-4" />
                    <span className="font-medium">{item.name}</span>
                  </div>
                  {item.badge && (
                    <span className="text-xs bg-emerald-500/20 text-emerald-400 px-2 py-1 rounded-full border border-emerald-500/30">
                      {item.badge}
                    </span>
                  )}
                </button>
              ))}
            </nav>
          </div>
        </div>

        <div className="lg:col-span-3 space-y-6">
          <div className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6">
            <h3 className="text-xl font-semibold text-white mb-6">Profile Settings</h3>
            <div className="space-y-6">
              <div className="flex items-center space-x-6">
                <div className="relative">
                  <div className="w-20 h-20 bg-gradient-to-br from-emerald-400 to-emerald-600 rounded-full flex items-center justify-center">
                    <span className="text-white font-bold text-xl">SC</span>
                  </div>
                  <button className="absolute -bottom-1 -right-1 w-6 h-6 bg-emerald-500 rounded-full flex items-center justify-center hover:bg-emerald-400 transition-colors">
                    <Plus className="w-3 h-3 text-white" />
                  </button>
                </div>
                <div className="flex-1">
                  <h4 className="text-white font-medium text-lg">Sarah Chen</h4>
                  <p className="text-zinc-400">sarah.chen@acme.com</p>
                  <p className="text-xs text-zinc-500 mt-1">Team Lead • San Francisco, CA • PST</p>
                  <div className="flex items-center space-x-4 mt-2">
                    <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-emerald-500/20 text-emerald-400 border border-emerald-500/30">
                      Pro Plan
                    </span>
                    <span className="text-xs text-zinc-400">Member since March 2023</span>
                  </div>
                </div>
                <button className="bg-zinc-800 hover:bg-zinc-700 text-white px-4 py-2 rounded-lg transition-colors">
                  Edit Profile
                </button>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-zinc-300 font-medium mb-2">First Name</label>
                  <input
                    type="text"
                    defaultValue="Sarah"
                    className="w-full px-4 py-3 bg-zinc-800/50 border border-zinc-700 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500 transition-colors"
                  />
                </div>
                <div>
                  <label className="block text-zinc-300 font-medium mb-2">Last Name</label>
                  <input
                    type="text"
                    defaultValue="Chen"
                    className="w-full px-4 py-3 bg-zinc-800/50 border border-zinc-700 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500 transition-colors"
                  />
                </div>
                <div>
                  <label className="block text-zinc-300 font-medium mb-2">Email</label>
                  <input
                    type="email"
                    defaultValue="sarah.chen@acme.com"
                    className="w-full px-4 py-3 bg-zinc-800/50 border border-zinc-700 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500 transition-colors"
                  />
                </div>
                <div>
                  <label className="block text-zinc-300 font-medium mb-2">Timezone</label>
                  <select className="w-full px-4 py-3 bg-zinc-800/50 border border-zinc-700 rounded-lg text-white focus:outline-none focus:border-emerald-500 transition-colors">
                    <option>Pacific Standard Time (PST)</option>
                    <option>Eastern Standard Time (EST)</option>
                    <option>Central Standard Time (CST)</option>
                    <option>Mountain Standard Time (MST)</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-zinc-300 font-medium mb-2">Bio</label>
                <textarea
                  rows={3}
                  defaultValue="Full-stack developer passionate about building scalable web applications. Team lead for the frontend infrastructure at Acme Corp."
                  className="w-full px-4 py-3 bg-zinc-800/50 border border-zinc-700 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500 transition-colors resize-none"
                />
              </div>
            </div>
          </div>

          <div className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold text-white">API Keys</h3>
              <button className="bg-emerald-500 hover:bg-emerald-400 text-white px-4 py-2 rounded-lg transition-all flex items-center space-x-2 text-sm">
                <Plus className="w-4 h-4" />
                <span>New Key</span>
              </button>
            </div>
            <div className="space-y-4">
              {[
                {
                  name: 'Production API Key',
                  key: 'kova_prod_ak_1a2b3c4d5e6f7g8h9i0j',
                  created: '2024-01-15',
                  lastUsed: '2 minutes ago',
                  permissions: ['deploy', 'read', 'write'],
                  requests: '2.4M this month'
                },
                {
                  name: 'Staging API Key',
                  key: 'kova_stag_ak_9z8y7x6w5v4u3t2s1r0q',
                  created: '2024-01-10',
                  lastUsed: '1 hour ago',
                  permissions: ['deploy', 'read'],
                  requests: '456K this month'
                },
                {
                  name: 'Development Key',
                  key: 'kova_dev_ak_p9o8i7u6y5t4r3e2w1q0',
                  created: '2024-01-05',
                  lastUsed: '3 days ago',
                  permissions: ['read'],
                  requests: '12.5K this month'
                },
              ].map((apiKey, index) => (
                <div key={index} className="p-4 bg-zinc-800/30 rounded-lg border border-zinc-700/50 hover:border-zinc-600/50 transition-colors">
                  <div className="flex items-center justify-between mb-3">
                    <div>
                      <h4 className="text-white font-medium">{apiKey.name}</h4>
                      <div className="flex items-center space-x-4 text-xs text-zinc-400 mt-1">
                        <span>Created {apiKey.created}</span>
                        <span>•</span>
                        <span>Last used {apiKey.lastUsed}</span>
                        <span>•</span>
                        <span>{apiKey.requests}</span>
                      </div>
                    </div>
                    <div className="flex items-center space-x-2">
                      <button
                        onClick={() => copyToClipboard(apiKey.key, `api-${index}`)}
                        className="text-zinc-400 hover:text-emerald-400 transition-colors"
                        title="Copy API Key"
                      >
                        {copied === `api-${index}` ?
                          <Check className="w-4 h-4 text-emerald-400" /> :
                          <Copy className="w-4 h-4" />
                        }
                      </button>
                      <button className="text-zinc-400 hover:text-white transition-colors" title="Edit">
                        <Eye className="w-4 h-4" />
                      </button>
                      <button className="text-zinc-400 hover:text-red-400 transition-colors" title="Delete">
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="font-mono text-sm text-zinc-300 bg-zinc-900/50 px-3 py-2 rounded border border-zinc-700/50">
                      {copied === `api-${index}` ? apiKey.key : `${apiKey.key.substring(0, 20)}••••••••••••••••••••`}
                    </div>
                    <div className="flex items-center space-x-2">
                      {apiKey.permissions.map((permission, idx) => (
                        <span key={idx} className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${permission === 'deploy' ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30' :
                          permission === 'write' ? 'bg-blue-500/20 text-blue-400 border border-blue-500/30' :
                            'bg-zinc-500/20 text-zinc-400 border border-zinc-500/30'
                          }`}>
                          {permission}
                        </span>
                      ))}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6">
              <h3 className="text-xl font-semibold text-white mb-6 flex items-center space-x-2">
                <Bell className="w-5 h-5 text-emerald-400" />
                <span>Notifications</span>
              </h3>
              <div className="space-y-4">
                {[
                  { label: 'Deployment Success', enabled: true, description: 'Get notified when deployments complete successfully' },
                  { label: 'Deployment Failures', enabled: true, description: 'Immediate alerts for failed deployments' },
                  { label: 'Build Warnings', enabled: false, description: 'Warnings during the build process' },
                  { label: 'Security Alerts', enabled: true, description: 'Security vulnerabilities in dependencies' },
                  { label: 'Weekly Reports', enabled: true, description: 'Weekly summary of deployment activity' },
                ].map((notification, index) => (
                  <div key={index} className="flex items-center justify-between p-3 bg-zinc-800/20 rounded-lg">
                    <div className="flex-1">
                      <div className="text-white font-medium text-sm">{notification.label}</div>
                      <div className="text-zinc-400 text-xs mt-1">{notification.description}</div>
                    </div>
                    <button className={`relative w-11 h-6 rounded-full transition-colors ${notification.enabled ? 'bg-emerald-500' : 'bg-zinc-600'
                      }`}>
                      <div className={`absolute top-1 w-4 h-4 bg-white rounded-full transition-transform ${notification.enabled ? 'translate-x-6' : 'translate-x-1'
                        }`}></div>
                    </button>
                  </div>
                ))}
              </div>
            </div>

            <div className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6">
              <h3 className="text-xl font-semibold text-white mb-6 flex items-center space-x-2">
                <Globe className="w-5 h-5 text-emerald-400" />
                <span>Integrations</span>
              </h3>
              <div className="space-y-4">
                {[
                  { name: 'GitHub', icon: Github, connected: true, status: 'Connected to 12 repositories' },
                  { name: 'GitLab', icon: Gitlab, connected: false, status: 'Connect to sync repositories' },
                  { name: 'Slack', icon: Bell, connected: true, status: 'Posting to #deployments channel' },
                  { name: 'Discord', icon: Bell, connected: false, status: 'Get deployment notifications' },
                  { name: 'Vercel', icon: Globe, connected: true, status: 'Import existing projects' },
                ].map((integration, index) => (
                  <div key={index} className="flex items-center justify-between p-3 bg-zinc-800/20 rounded-lg">
                    <div className="flex items-center space-x-3">
                      <integration.icon className="w-5 h-5 text-zinc-400" />
                      <div>
                        <div className="text-white font-medium text-sm">{integration.name}</div>
                        <div className="text-zinc-400 text-xs">{integration.status}</div>
                      </div>
                    </div>
                    <button className={`px-3 py-1 rounded-lg text-xs font-medium transition-colors ${integration.connected
                      ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 hover:bg-emerald-500/30'
                      : 'bg-zinc-700 text-white hover:bg-zinc-600'
                      }`}>
                      {integration.connected ? 'Connected' : 'Connect'}
                    </button>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <div className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6">
            <h3 className="text-xl font-semibold text-white mb-6 flex items-center space-x-2">
              <Users className="w-5 h-5 text-emerald-400" />
              <span>Team Management</span>
            </h3>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="text-white font-medium">Team Members</h4>
                  <p className="text-zinc-400 text-sm">Manage access and permissions for your team</p>
                </div>
                <button className="bg-emerald-500 hover:bg-emerald-400 text-white px-4 py-2 rounded-lg transition-all flex items-center space-x-2 text-sm">
                  <Plus className="w-4 h-4" />
                  <span>Invite Member</span>
                </button>
              </div>

              <div className="space-y-3">
                {[
                  { name: 'Sarah Chen', email: 'sarah.chen@acme.com', role: 'Owner', avatar: 'SC', status: 'active' },
                  { name: 'Mike Johnson', email: 'mike.johnson@acme.com', role: 'Admin', avatar: 'MJ', status: 'active' },
                  { name: 'Emily Rodriguez', email: 'emily.rodriguez@acme.com', role: 'Developer', avatar: 'ER', status: 'active' },
                  { name: 'David Kim', email: 'david.kim@acme.com', role: 'Developer', avatar: 'DK', status: 'pending' },
                ].map((member, index) => (
                  <div key={index} className="flex items-center justify-between p-3 bg-zinc-800/20 rounded-lg">
                    <div className="flex items-center space-x-3">
                      <div className="w-10 h-10 bg-gradient-to-br from-emerald-400 to-emerald-600 rounded-full flex items-center justify-center">
                        <span className="text-white font-medium text-sm">{member.avatar}</span>
                      </div>
                      <div>
                        <div className="text-white font-medium text-sm">{member.name}</div>
                        <div className="text-zinc-400 text-xs">{member.email}</div>
                      </div>
                    </div>
                    <div className="flex items-center space-x-3">
                      <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${member.status === 'active'
                        ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
                        : 'bg-yellow-500/20 text-yellow-400 border border-yellow-500/30'
                        }`}>
                        {member.status}
                      </span>
                      <select className="bg-zinc-700 text-white text-xs px-2 py-1 rounded focus:outline-none focus:border-emerald-500">
                        <option value="owner">Owner</option>
                        <option value="admin">Admin</option>
                        <option value="developer">Developer</option>
                        <option value="viewer">Viewer</option>
                      </select>
                      <button className="text-zinc-400 hover:text-red-400 transition-colors">
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <div className="bg-zinc-900/50 backdrop-blur-sm border border-zinc-800/50 rounded-xl p-6">
            <h3 className="text-xl font-semibold text-white mb-6">Danger Zone</h3>
            <div className="space-y-4">
              <div className="border border-red-800/50 rounded-lg p-4 bg-red-950/20">
                <div className="flex items-start space-x-4">
                  <AlertTriangle className="w-5 h-5 text-red-400 mt-0.5" />
                  <div className="flex-1">
                    <h4 className="text-red-400 font-medium mb-2">Delete Account</h4>
                    <p className="text-zinc-400 text-sm mb-4">
                      Permanently delete your account and all associated data. This action cannot be undone and will:
                    </p>
                    <ul className="text-zinc-400 text-sm space-y-1 mb-4 ml-4">
                      <li>• Delete all projects and deployments</li>
                      <li>• Remove all team members and API keys</li>
                      <li>• Cancel your subscription and stop all billing</li>
                      <li>• Permanently delete all data within 30 days</li>
                    </ul>
                    <button className="bg-red-600 hover:bg-red-500 text-white px-4 py-2 rounded-lg transition-colors text-sm">
                      Delete Account
                    </button>
                  </div>
                </div>
              </div>

              <div className="border border-orange-800/50 rounded-lg p-4 bg-orange-950/20">
                <div className="flex items-start space-x-4">
                  <AlertCircle className="w-5 h-5 text-orange-400 mt-0.5" />
                  <div className="flex-1">
                    <h4 className="text-orange-400 font-medium mb-2">Reset All API Keys</h4>
                    <p className="text-zinc-400 text-sm mb-4">
                      Regenerate all API keys for security purposes. All existing keys will be invalidated immediately.
                    </p>
                    <button className="bg-orange-600 hover:bg-orange-500 text-white px-4 py-2 rounded-lg transition-colors text-sm">
                      Reset All Keys
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );

  const renderContent = () => {
    switch (activeTab) {
      case 'dashboard':
        return <DashboardPage />;
      case 'projects':
        return <ProjectsPage />;
      case 'deployments':
        return <DeploymentsPage />;
      case 'settings':
        return <SettingsPage />;
      default:
        return <DashboardPage />;
    }
  };

  return (
    <div className="min-h-screen bg-zinc-950 relative overflow-hidden">
      {/* Background Effects */}
      <div className="absolute inset-0">
        <div className="absolute inset-0 bg-gradient-to-br from-slate-950/30 via-transparent to-emerald-950/15"></div>
        <div
          className="absolute inset-0 opacity-5"
          style={{
            backgroundImage: `radial-gradient(circle at 50% 50%, #10b981 1px, transparent 1px)`,
            backgroundSize: '100px 100px'
          }}
        ></div>
        <div className="absolute top-1/3 right-1/4 w-96 h-96 bg-emerald-500/3 rounded-full blur-3xl"></div>
        <div className="absolute bottom-1/4 left-1/3 w-64 h-64 bg-blue-500/2 rounded-full blur-2xl"></div>
      </div>

      <div className="relative z-10 flex">
        <Sidebar />
        <main className="flex-1 p-8 max-w-full overflow-x-auto">
          {renderContent()}
        </main>
      </div>
    </div>
  );
};

export default KovaDashboard;
