"use client"
import React, { useState, useEffect } from 'react';
import { NavbarProps, NavigationState } from './types';
import { DesktopNavbar } from './DesktopNavbar';
import { MobileNavbar } from './MobileNavbar';
import { TabletNavbar } from './TabletNavbar';
import { CommandPalette } from './CommandPalette';
import { useKeyboardShortcuts } from '@/hooks/navbar';
import { useBreakpoint } from '@/hooks/navbar';
import { useSession } from 'next-auth/react';

export const Navbar: React.FC<NavbarProps> = ({
  currentProject,
  projects = [],
  onNavigate,
  onProjectChange,
  onNewProject,
  className = '',
  notifications = 0,
  systemStatus = 'operational'
}) => {
  const { data } = useSession()
  const user = {
    email: data?.user.email || "t@t.com",
    name: data?.user.username || "test",
    id: data?.user.id || "1",
    role: "admin"

  }
  const [navigationState, setNavigationState] = useState<NavigationState>({
    activeTab: 'dashboard',
    activeSubTab: '',
    showCommand: false,
    showProjectSwitcher: false,
    showUserMenu: false,
    showMobileMenu: false,
    searchQuery: '',
    viewMode: 'table',
    favorites: [currentProject].filter(Boolean) as string[]
  });

  const { isMobile, isTablet, isDesktop } = useBreakpoint();

  // Update navigation state
  const updateState = (updates: Partial<NavigationState>) => {
    setNavigationState(prev => ({ ...prev, ...updates }));
  };

  // Handle navigation changes
  const handleNavigate = (tab: string, subTab?: string) => {
    updateState({
      activeTab: tab,
      activeSubTab: subTab || '',
      showMobileMenu: false
    });
    onNavigate?.(tab, subTab);
  };

  // Keyboard shortcuts
  useKeyboardShortcuts({
    navigationState,
    updateState,
    handleNavigate,
    onProjectChange,
    onNewProject,
    currentProject
  });

  // Close menus on outside click
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as Element;
      if (!target.closest('[data-menu]')) {
        updateState({
          showUserMenu: false,
          showProjectSwitcher: false
        });
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const sharedProps = {
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
  };

  return (
    <div className={className}>
      {/* Render appropriate navbar based on screen size */}
      {isMobile && <MobileNavbar {...sharedProps} />}
      {isTablet && <TabletNavbar {...sharedProps} />}
      {isDesktop && <DesktopNavbar {...sharedProps} />}

      {/* Command Palette (shared across all sizes) */}
      <CommandPalette
        navigationState={navigationState}
        updateState={updateState}
        handleNavigate={handleNavigate}
        currentProject={currentProject}
        projects={projects}
      />
    </div>
  );
};

export default Navbar;
export type { NavbarProps, NavigationState } from './types';
