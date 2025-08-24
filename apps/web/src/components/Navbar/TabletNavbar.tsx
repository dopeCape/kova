import React from 'react';
import {
  Code, Search, Plus, Bell, User, ChevronDown, Menu, X,
  Filter, MoreVertical, Folder, ChevronRight, Star
} from 'lucide-react';
import { SharedNavbarProps } from './types';
import { mainNavItems, getSubNavigation, getBreadcrumbs } from './navigationConfig';

export const TabletNavbar: React.FC<SharedNavbarProps> = ({
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

  return (
    <>
      {/* Tablet Header */}
      <header className="sticky top-0 z-50 bg-zinc-950/95 backdrop-blur-sm border-b border-zinc-800">
        <div className="flex items-center justify-between h-16 px-4">

          {/* Logo + Menu Toggle */}
          <div className="flex items-center space-x-4">
            <button
              onClick={() => updateState({ showMobileMenu: !navigationState.showMobileMenu })}
              className="p-2 text-zinc-400 hover:text-white transition-colors md:hidden"
            >
              {navigationState.showMobileMenu ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </button>

            <div className="flex items-center space-x-3">
              <div className="w-7 h-7 bg-emerald-500 rounded-md flex items-center justify-center">
                <Code className="w-4 h-4 text-white" />
              </div>
              <div className="font-bold text-lg">Kova</div>
            </div>

            {/* Compact Navigation for larger tablets */}
            <nav className="hidden md:flex items-center space-x-1">
              {mainNavItems.slice(0, 4).map((item) => (
                <button
                  key={item.id}
                  onClick={() => handleNavigate(item.id)}
                  className={`flex items-center space-x-2 px-3 py-2 text-sm font-medium transition-colors relative ${navigationState.activeTab === item.id
                    ? 'text-white'
                    : 'text-zinc-400 hover:text-white'
                    }`}
                >
                  <item.icon className="w-4 h-4" />
                  <span className="hidden lg:block">{item.name}</span>

                  {navigationState.activeTab === item.id && (
                    <div className="absolute bottom-0 left-3 right-3 h-0.5 bg-emerald-500 rounded-full" />
                  )}
                </button>
              ))}
            </nav>
          </div>

          {/* Right Actions */}
          <div className="flex items-center space-x-3">

            {/* Search */}
            <button
              onClick={() => updateState({ showCommand: true })}
              className="p-2 text-zinc-400 hover:text-white transition-colors"
            >
              <Search className="w-5 h-5" />
            </button>

            {/* New Project */}
            <button
              onClick={onNewProject}
              className="flex items-center space-x-2 bg-emerald-600 hover:bg-emerald-500 text-white px-3 py-2 rounded-lg transition-colors text-sm font-medium"
            >
              <Plus className="w-4 h-4" />
              <span className="hidden sm:block">New</span>
            </button>

            {/* Status */}
            <div className={`w-2 h-2 rounded-full ${systemStatus === 'operational' ? 'bg-emerald-500' :
              systemStatus === 'degraded' ? 'bg-yellow-500' : 'bg-red-500'
              }`} />

            {/* Notifications */}
            <button className="p-2 text-zinc-400 hover:text-white transition-colors relative">
              <Bell className="w-5 h-5" />
              {notifications > 0 && (
                <div className="absolute -top-1 -right-1 w-4 h-4 bg-emerald-500 rounded-full flex items-center justify-center text-xs text-white font-medium">
                  {notifications > 9 ? '9+' : notifications}
                </div>
              )}
            </button>

            {/* User */}
            <button
              onClick={() => updateState({ showUserMenu: !navigationState.showUserMenu })}
              className="flex items-center space-x-2 text-sm hover:bg-zinc-900 p-2 rounded-lg transition-colors"
            >
              <div className="w-6 h-6 bg-zinc-700 rounded-full flex items-center justify-center">
                <User className="w-3 h-3" />
              </div>
              <span className="text-zinc-300 hidden sm:block">{user.name.split(' ')[0]}</span>
            </button>
          </div>
        </div>

        {/* Sub-Navigation (Horizontal Scroll) */}
        {subNavigation.length > 0 && (
          <div className="border-t border-zinc-800 bg-zinc-900/20">
            <div className="px-4">
              <div className="flex items-center space-x-1 overflow-x-auto scrollbar-hide py-2">
                {subNavigation.map((item) => (
                  <button
                    key={item.id}
                    onClick={() => updateState({ activeSubTab: item.id })}
                    className={`flex items-center space-x-2 px-3 py-2 text-sm font-medium transition-colors whitespace-nowrap ${navigationState.activeSubTab === item.id
                      ? 'text-white bg-zinc-800 rounded-lg'
                      : 'text-zinc-400 hover:text-white'
                      }`}
                  >
                    <item.icon className="w-4 h-4" />
                    <span>{item.name}</span>
                    {item.badge && (
                      <span className={`text-xs px-1.5 py-0.5 rounded-full ${item.badgeColor === 'emerald' ? 'bg-emerald-500/20 text-emerald-400' :
                        item.badgeColor === 'yellow' ? 'bg-yellow-500/20 text-yellow-400' :
                          item.badgeColor === 'red' ? 'bg-red-500/20 text-red-400' :
                            'bg-zinc-700 text-zinc-300'
                        }`}>
                        {item.badge}
                      </span>
                    )}
                  </button>
                ))}
              </div>
            </div>
          </div>
        )}

        {/* Breadcrumbs (Minimal) */}
        {breadcrumbs.length > 1 && (
          <div className="border-t border-zinc-800 bg-zinc-900/10">
            <div className="px-4 py-2">
              <div className="flex items-center space-x-2 text-sm overflow-x-auto scrollbar-hide">
                {breadcrumbs.slice(-2).map((crumb, index) => (
                  <div key={crumb.path} className="flex items-center space-x-2 whitespace-nowrap">
                    {index > 0 && <ChevronRight className="w-3 h-3 text-zinc-600" />}
                    <button
                      className={`flex items-center space-x-1 px-2 py-1 rounded transition-colors ${crumb.current
                        ? 'text-emerald-300 bg-emerald-500/10'
                        : 'text-zinc-400 hover:text-white'
                        }`}
                    >
                      <crumb.icon className="w-3 h-3" />
                      <span>{crumb.name}</span>
                    </button>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </header>

      {/* Full Screen Menu for Mobile Screens */}
      {navigationState.showMobileMenu && (
        <div className="fixed inset-0 z-50 bg-zinc-950 md:hidden">
          <div className="flex flex-col h-full">

            {/* Header */}
            <div className="flex items-center justify-between p-4 border-b border-zinc-800">
              <div className="flex items-center space-x-3">
                <div className="w-8 h-8 bg-emerald-500 rounded-lg flex items-center justify-center">
                  <Code className="w-4 h-4 text-white" />
                </div>
                <div>
                  <div className="font-bold text-lg">Kova</div>
                  <div className="text-xs text-zinc-400">Deploy Platform</div>
                </div>
              </div>
              <button
                onClick={() => updateState({ showMobileMenu: false })}
                className="p-2 text-zinc-400"
              >
                <X className="w-6 h-6" />
              </button>
            </div>

            {/* Content */}
            <div className="flex-1 overflow-y-auto p-4">

              {/* User Info */}
              <div className="mb-6 p-4 bg-zinc-900 rounded-lg border border-zinc-800">
                <div className="flex items-center space-x-3">
                  <div className="w-12 h-12 bg-zinc-700 rounded-full flex items-center justify-center">
                    <User className="w-6 h-6" />
                  </div>
                  <div>
                    <div className="font-medium text-white">{user.name}</div>
                    <div className="text-sm text-zinc-400">{user.email}</div>
                    <div className="text-xs text-zinc-500">{user.role}</div>
                  </div>
                </div>
              </div>

              {/* Main Navigation */}
              <div className="space-y-1 mb-6">
                {mainNavItems.map((item) => (
                  <button
                    key={item.id}
                    onClick={() => {
                      handleNavigate(item.id);
                      updateState({ showMobileMenu: false });
                    }}
                    className={`w-full flex items-center space-x-3 p-4 rounded-lg transition-colors ${navigationState.activeTab === item.id
                      ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30'
                      : 'text-zinc-300 hover:bg-zinc-800'
                      }`}
                  >
                    <item.icon className="w-5 h-5" />
                    <span className="font-medium">{item.name}</span>
                  </button>
                ))}
              </div>

              {/* Current Project (if applicable) */}
              {currentProject && navigationState.activeTab === 'projects' && (
                <div className="mb-6">
                  <div className="text-xs text-zinc-400 uppercase tracking-wide px-2 py-2">Current Project</div>
                  <div className="p-4 bg-zinc-900 rounded-lg border border-zinc-800">
                    <div className="flex items-center space-x-3">
                      <Folder className="w-5 h-5 text-emerald-400" />
                      <div>
                        <div className="font-mono text-white">{currentProject}</div>
                        <div className="text-xs text-zinc-400">Tap to switch projects</div>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>

            {/* Footer Actions */}
            <div className="border-t border-zinc-800 p-4 space-y-2">
              <button className="w-full p-3 text-center text-sm text-zinc-400 hover:text-white hover:bg-zinc-900 rounded-lg transition-colors">
                Account Settings
              </button>
              <button className="w-full p-3 text-center text-sm text-zinc-400 hover:text-white hover:bg-zinc-900 rounded-lg transition-colors">
                Sign Out
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Custom scrollbar styles */}
      <style jsx>{`
        .scrollbar-hide {
          -ms-overflow-style: none;
          scrollbar-width: none;
        }
        .scrollbar-hide::-webkit-scrollbar {
          display: none;
        }
      `}</style>
    </>
  );
};
