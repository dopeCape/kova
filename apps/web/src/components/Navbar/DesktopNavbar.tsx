import React from 'react';
import {
  Code, Search, Plus, Bell, User, ChevronDown, Filter,
  Share, MoreVertical, Folder, ChevronRight, Star,
  List, Grid3X3
} from 'lucide-react';
import { SharedNavbarProps } from './types';
import { mainNavItems, getSubNavigation, getBreadcrumbs } from './navigationConfig';

export const DesktopNavbar: React.FC<SharedNavbarProps> = ({
  navigationState,
  updateState,
  handleNavigate,
  user,
  currentProject,
  projects,
  onProjectChange,
  onNewProject,
  notifications,
  systemStatus
}) => {
  const subNavigation = getSubNavigation(navigationState.activeTab, currentProject);
  const breadcrumbs = getBreadcrumbs(navigationState.activeTab, navigationState.activeSubTab, currentProject);
  const showDataViewControls = ['projects', 'deployments', 'team'].includes(navigationState.activeTab);

  return (
    <header className="sticky top-0 z-50 bg-zinc-950/95 backdrop-blur-sm border-b border-zinc-800">
      {/* Main Navigation Bar */}
      <div className="flex items-center justify-between h-16 px-6 max-w-[2000px] mx-auto">

        {/* Logo + Main Navigation */}
        <div className="flex items-center space-x-8">
          <div className="flex items-center space-x-3">
            <div className="w-7 h-7 bg-emerald-500 rounded-md flex items-center justify-center">
              <Code className="w-4 h-4 text-white" />
            </div>
            <div className="font-bold text-lg">Kova</div>
          </div>

          {/* Main Navigation Tabs */}
          <nav className="flex items-center space-x-1">
            {mainNavItems.map((item) => (
              <button
                key={item.id}
                onClick={() => handleNavigate(item.id)}
                className={`flex items-center space-x-2 px-4 py-2 text-sm font-medium transition-colors relative ${navigationState.activeTab === item.id
                  ? 'text-white'
                  : 'text-zinc-400 hover:text-white'
                  }`}
              >
                <item.icon className="w-4 h-4" />
                <span>{item.name}</span>
                <span className="text-xs text-zinc-500">⌘{item.shortcut}</span>

                {navigationState.activeTab === item.id && (
                  <div className="absolute bottom-0 left-4 right-4 h-0.5 bg-emerald-500 rounded-full" />
                )}
              </button>
            ))}
          </nav>
        </div>

        {/* Right Side Actions */}
        <div className="flex items-center space-x-4">

          {/* Command Search */}
          <button
            onClick={() => updateState({ showCommand: true })}
            className="flex items-center space-x-3 px-4 py-2 bg-zinc-900 hover:bg-zinc-800 border border-zinc-700 hover:border-zinc-600 rounded-lg transition-all text-sm text-zinc-400"
          >
            <Search className="w-4 h-4" />
            <span>Search...</span>
            <div className="flex items-center space-x-1">
              <kbd className="px-1.5 py-0.5 bg-zinc-800 border border-zinc-600 rounded text-xs">⌘</kbd>
              <kbd className="px-1.5 py-0.5 bg-zinc-800 border border-zinc-600 rounded text-xs">K</kbd>
            </div>
          </button>

          {/* View Mode Switcher */}
          {showDataViewControls && (
            <div className="flex items-center bg-zinc-900 border border-zinc-700 rounded-lg p-1">
              <button
                onClick={() => updateState({ viewMode: 'table' })}
                className={`p-1.5 rounded transition-colors ${navigationState.viewMode === 'table' ? 'bg-zinc-700 text-white' : 'text-zinc-400 hover:text-white'
                  }`}
              >
                <List className="w-4 h-4" />
              </button>
              <button
                onClick={() => updateState({ viewMode: 'grid' })}
                className={`p-1.5 rounded transition-colors ${navigationState.viewMode === 'grid' ? 'bg-zinc-700 text-white' : 'text-zinc-400 hover:text-white'
                  }`}
              >
                <Grid3X3 className="w-4 h-4" />
              </button>
            </div>
          )}

          {/* New Button - Minimal */}
          <button
            onClick={onNewProject}
            className="flex items-center space-x-2 bg-emerald-600 hover:bg-emerald-500 text-white px-3 py-2 rounded-lg transition-colors text-sm font-medium"
          >
            <Plus className="w-4 h-4" />
            <span>New</span>
          </button>

          {/* Status + User */}
          <div className="flex items-center space-x-4 border-l border-zinc-700 pl-4">
            <div className="flex items-center space-x-2 text-sm text-zinc-400">
              <div className={`w-2 h-2 rounded-full animate-pulse ${systemStatus === 'operational' ? 'bg-emerald-500' :
                systemStatus === 'degraded' ? 'bg-yellow-500' : 'bg-red-500'
                }`} />
              <span className="capitalize">{systemStatus}</span>
            </div>

            <button className="p-2 text-zinc-400 hover:text-white transition-colors relative">
              <Bell className="w-5 h-5" />
              {notifications > 0 && (
                <div className="absolute -top-1 -right-1 w-5 h-5 bg-emerald-500 rounded-full flex items-center justify-center text-xs text-white font-medium">
                  {notifications > 9 ? '9+' : notifications}
                </div>
              )}
            </button>

            <div className="relative" data-menu>
              <button
                onClick={() => updateState({ showUserMenu: !navigationState.showUserMenu })}
                className="flex items-center space-x-2 text-sm hover:bg-zinc-900 p-2 rounded-lg transition-colors"
              >
                <div className="w-7 h-7 bg-zinc-700 rounded-full flex items-center justify-center">
                  <User className="w-4 h-4" />
                </div>
                <span className="text-zinc-300">{user.name.split(' ')[0]}</span>
                <ChevronDown className="w-3 h-3 text-zinc-400" />
              </button>

              {navigationState.showUserMenu && (
                <div className="absolute right-0 top-full mt-2 w-64 bg-zinc-900 border border-zinc-700 rounded-lg shadow-xl py-2 z-50">
                  <div className="px-4 py-3 border-b border-zinc-700">
                    <div className="font-medium text-white">{user.name}</div>
                    <div className="text-sm text-zinc-400">{user.email}</div>
                    <div className="text-xs text-zinc-500 mt-1">{user.role}</div>
                  </div>
                  <div className="py-2">
                    <button className="w-full px-4 py-2 text-left text-sm text-zinc-300 hover:bg-zinc-800 transition-colors">
                      Account Settings
                    </button>
                    <button className="w-full px-4 py-2 text-left text-sm text-zinc-300 hover:bg-zinc-800 transition-colors">
                      Billing & Usage
                    </button>
                    <button className="w-full px-4 py-2 text-left text-sm text-zinc-300 hover:bg-zinc-800 transition-colors">
                      Team Settings
                    </button>
                  </div>
                  <div className="border-t border-zinc-700 pt-2">
                    <button className="w-full px-4 py-2 text-left text-sm text-zinc-400 hover:bg-zinc-800 transition-colors">
                      Sign Out
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Sub-Navigation (ONLY ONE - no duplicates) */}
      {subNavigation.length > 0 && (
        <div className="border-t border-zinc-800 bg-zinc-900/20">
          <div className="flex items-center justify-between px-6 max-w-[2000px] mx-auto">
            <nav className="flex items-center">
              {subNavigation.map((item) => (
                <button
                  key={item.id}
                  onClick={() => updateState({ activeSubTab: item.id })}
                  className={`flex items-center space-x-2 px-4 py-3 text-sm font-medium transition-colors relative ${navigationState.activeSubTab === item.id
                    ? 'text-white border-b-2 border-emerald-500'
                    : 'text-zinc-400 hover:text-white'
                    }`}
                >
                  <item.icon className="w-4 h-4" />
                  <span>{item.name}</span>
                  {item.badge && (
                    <span className={`text-xs px-2 py-0.5 rounded-full ${item.badgeColor === 'emerald' ? 'bg-emerald-500/20 text-emerald-400' :
                      item.badgeColor === 'yellow' ? 'bg-yellow-500/20 text-yellow-400' :
                        item.badgeColor === 'red' ? 'bg-red-500/20 text-red-400' :
                          'bg-zinc-700 text-zinc-300'
                      }`}>
                      {item.badge}
                    </span>
                  )}
                  <span className="text-xs text-zinc-500">⌥{item.shortcut}</span>
                </button>
              ))}
            </nav>

            {/* Sub-navigation Actions */}
            <div className="flex items-center space-x-2">
              <button className="p-2 text-zinc-400 hover:text-white transition-colors">
                <Filter className="w-4 h-4" />
              </button>
              <button className="p-2 text-zinc-400 hover:text-white transition-colors">
                <Share className="w-4 h-4" />
              </button>
              <button className="p-2 text-zinc-400 hover:text-white transition-colors">
                <MoreVertical className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Breadcrumbs - After Sub-Navigation (ONLY ONE - no duplicates) */}
      {breadcrumbs.length > 1 && (
        <div className="border-t border-zinc-800 bg-zinc-900/10">
          <div className="flex items-center justify-between px-6 py-2 max-w-[2000px] mx-auto">

            {/* Breadcrumb Navigation */}
            <nav className="flex items-center space-x-2 text-sm">
              {breadcrumbs.map((crumb, index) => (
                <div key={crumb.path} className="flex items-center space-x-2">
                  {index > 0 && <ChevronRight className="w-3 h-3 text-zinc-600" />}
                  <button
                    className={`flex items-center space-x-1 px-2 py-1 rounded transition-colors ${crumb.current
                      ? 'text-emerald-300 bg-emerald-500/10'
                      : 'text-zinc-400 hover:text-white hover:bg-zinc-900'
                      }`}
                  >
                    <crumb.icon className="w-3 h-3" />
                    <span>{crumb.name}</span>
                  </button>
                </div>
              ))}
            </nav>

            {/* Project Switcher (only when in project context) */}
            {navigationState.activeTab === 'projects' && currentProject && (
              <div className="relative" data-menu>
                <button
                  onClick={() => updateState({ showProjectSwitcher: !navigationState.showProjectSwitcher })}
                  className="flex items-center space-x-2 px-3 py-1.5 bg-zinc-900 hover:bg-zinc-800 border border-zinc-700 rounded-lg transition-colors text-sm"
                >
                  <Folder className="w-4 h-4 text-emerald-400" />
                  <span className="font-mono text-emerald-300">{currentProject}</span>
                  <ChevronDown className="w-3 h-3 text-zinc-400" />
                  <div className="flex items-center space-x-1 border-l border-zinc-700 pl-2 ml-2">
                    <kbd className="px-1 py-0.5 bg-zinc-800 border border-zinc-600 rounded text-xs text-zinc-500">⌘P</kbd>
                  </div>
                </button>

                {navigationState.showProjectSwitcher && (
                  <div className="absolute right-0 top-full mt-2 w-80 bg-zinc-900 border border-zinc-700 rounded-lg shadow-xl z-50">
                    <div className="p-4 border-b border-zinc-700">
                      <div className="relative">
                        <Search className="w-4 h-4 absolute left-3 top-1/2 transform -translate-y-1/2 text-zinc-400" />
                        <input
                          type="text"
                          placeholder="Search projects..."
                          value={navigationState.searchQuery}
                          onChange={(e) => updateState({ searchQuery: e.target.value })}
                          className="w-full pl-10 pr-4 py-2 bg-zinc-800 border border-zinc-700 rounded text-white placeholder-zinc-400 focus:outline-none focus:border-emerald-500"
                        />
                      </div>
                    </div>

                    <div className="max-h-64 overflow-y-auto p-2">
                      {projects
                        .filter(project =>
                          project.name.toLowerCase().includes(navigationState.searchQuery.toLowerCase()) ||
                          project.team.toLowerCase().includes(navigationState.searchQuery.toLowerCase())
                        )
                        .map((project) => (
                          <button
                            key={project.id}
                            onClick={() => {
                              onProjectChange?.(project.id);
                              updateState({ showProjectSwitcher: false, searchQuery: '' });
                            }}
                            className="w-full flex items-center space-x-3 p-3 hover:bg-zinc-800 rounded transition-colors group"
                          >
                            <div className={`w-2 h-2 rounded-full ${project.status === 'active' ? 'bg-emerald-500' :
                              project.status === 'error' ? 'bg-red-500' :
                                project.status === 'building' ? 'bg-yellow-500 animate-pulse' : 'bg-zinc-500'
                              }`} />
                            <div className="flex-1 text-left">
                              <div className="text-white font-mono text-sm">{project.name}</div>
                              <div className="text-zinc-400 text-xs">{project.team} • {project.lastDeploy}</div>
                            </div>
                            <button
                              onClick={(e) => {
                                e.stopPropagation();
                                const newFavorites = navigationState.favorites.includes(project.id)
                                  ? navigationState.favorites.filter(id => id !== project.id)
                                  : [...navigationState.favorites, project.id];
                                updateState({ favorites: newFavorites });
                              }}
                              className="opacity-0 group-hover:opacity-100 transition-opacity"
                            >
                              <Star className={`w-4 h-4 ${navigationState.favorites.includes(project.id)
                                ? 'text-yellow-500 fill-current'
                                : 'text-zinc-400'
                                }`} />
                            </button>
                          </button>
                        ))}
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      )}
    </header>
  );
};
