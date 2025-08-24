import React from 'react';
import {
  Code, Menu, X, Bell, User, Plus, Search,
  Activity, Rocket, Settings, BarChart3
} from 'lucide-react';
import { SharedNavbarProps } from './types';
import { mainNavItems } from './navigationConfig';

export const MobileNavbar: React.FC<SharedNavbarProps> = ({
  navigationState,
  updateState,
  handleNavigate,
  user,
  onNewProject,
  notifications,
  systemStatus
}) => {
  // Simplified navigation for mobile - only show essential items
  const mobileNavItems = mainNavItems.slice(0, 4); // Dashboard, Projects, Deployments, Settings

  return (
    <>
      {/* Mobile Header */}
      <header className="sticky top-0 z-50 bg-zinc-950/95 backdrop-blur-sm border-b border-zinc-800">
        <div className="flex items-center justify-between h-14 px-4">

          {/* Logo */}
          <div className="flex items-center space-x-3">
            <div className="w-6 h-6 bg-emerald-500 rounded-md flex items-center justify-center">
              <Code className="w-3 h-3 text-white" />
            </div>
            <div className="font-bold">Kova</div>
          </div>

          {/* Right Actions */}
          <div className="flex items-center space-x-3">
            {/* Status Indicator */}
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

            {/* Menu Toggle */}
            <button
              onClick={() => updateState({ showMobileMenu: !navigationState.showMobileMenu })}
              className="p-2 text-zinc-400 hover:text-white transition-colors"
            >
              {navigationState.showMobileMenu ? (
                <X className="w-5 h-5" />
              ) : (
                <Menu className="w-5 h-5" />
              )}
            </button>
          </div>
        </div>
      </header>

      {/* Mobile Menu Overlay */}
      {navigationState.showMobileMenu && (
        <div className="fixed inset-0 z-50 bg-zinc-950">
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

            {/* Navigation */}
            <div className="flex-1 overflow-y-auto p-4">

              {/* User Info */}
              <div className="mb-6 p-4 bg-zinc-900 rounded-lg border border-zinc-800">
                <div className="flex items-center space-x-3">
                  <div className="w-10 h-10 bg-zinc-700 rounded-full flex items-center justify-center">
                    <User className="w-5 h-5" />
                  </div>
                  <div>
                    <div className="font-medium text-white">{user.name}</div>
                    <div className="text-sm text-zinc-400">{user.role}</div>
                  </div>
                </div>
              </div>

              {/* Quick Actions */}
              <div className="mb-6 space-y-2">
                <button
                  onClick={() => {
                    onNewProject?.();
                    updateState({ showMobileMenu: false });
                  }}
                  className="w-full flex items-center space-x-3 p-4 bg-emerald-600 hover:bg-emerald-500 text-white rounded-lg transition-colors"
                >
                  <Plus className="w-5 h-5" />
                  <span className="font-medium">New Project</span>
                </button>

                <button
                  onClick={() => {
                    updateState({ showCommand: true, showMobileMenu: false });
                  }}
                  className="w-full flex items-center space-x-3 p-4 bg-zinc-900 hover:bg-zinc-800 border border-zinc-700 text-white rounded-lg transition-colors"
                >
                  <Search className="w-5 h-5" />
                  <span className="font-medium">Search</span>
                </button>
              </div>

              {/* Main Navigation */}
              <div className="space-y-1">
                <div className="text-xs text-zinc-400 uppercase tracking-wide px-2 py-2">Navigation</div>
                {mobileNavItems.map((item) => (
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

              {/* Quick Links */}
              <div className="mt-6 pt-6 border-t border-zinc-800">
                <div className="text-xs text-zinc-400 uppercase tracking-wide px-2 py-2">Quick Links</div>
                <div className="space-y-1">
                  <button className="w-full flex items-center space-x-3 p-3 text-zinc-300 hover:bg-zinc-800 rounded-lg transition-colors text-left">
                    <span>Account Settings</span>
                  </button>
                  <button className="w-full flex items-center space-x-3 p-3 text-zinc-300 hover:bg-zinc-800 rounded-lg transition-colors text-left">
                    <span>Team Settings</span>
                  </button>
                  <button className="w-full flex items-center space-x-3 p-3 text-zinc-300 hover:bg-zinc-800 rounded-lg transition-colors text-left">
                    <span>Billing</span>
                  </button>
                </div>
              </div>
            </div>

            {/* Footer */}
            <div className="border-t border-zinc-800 p-4">
              <button className="w-full text-center text-sm text-zinc-400 hover:text-white transition-colors">
                Sign Out
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};
