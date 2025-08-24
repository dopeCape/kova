import { useEffect, useState } from 'react';
import { NavigationState } from '@/components/Navbar/types';
import { mainNavItems, getSubNavigation } from '@/components/Navbar/navigationConfig';

// Keyboard shortcuts hook
export const useKeyboardShortcuts = ({
  navigationState,
  updateState,
  handleNavigate,
  onProjectChange,
  onNewProject,
  currentProject
}: {
  navigationState: NavigationState;
  updateState: (updates: Partial<NavigationState>) => void;
  handleNavigate: (tab: string, subTab?: string) => void;
  onProjectChange?: (projectId: string) => void;
  onNewProject?: () => void;
  currentProject?: string;
}) => {
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Command palette
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        updateState({ showCommand: !navigationState.showCommand });
      }

      // Project switcher
      if ((e.metaKey || e.ctrlKey) && e.key === 'p') {
        e.preventDefault();
        updateState({ showProjectSwitcher: !navigationState.showProjectSwitcher });
      }

      // Close modals
      if (e.key === 'Escape') {
        updateState({
          showCommand: false,
          showProjectSwitcher: false,
          showUserMenu: false,
          showMobileMenu: false
        });
      }

      // Main navigation shortcuts (⌘1-7)
      if ((e.metaKey || e.ctrlKey) && ['1', '2', '3', '4', '5', '6', '7'].includes(e.key)) {
        e.preventDefault();
        const navItem = mainNavItems[parseInt(e.key) - 1];
        if (navItem) {
          handleNavigate(navItem.id);
        }
      }

      // Sub navigation shortcuts (⌥1-5)
      if (e.altKey && ['1', '2', '3', '4', '5'].includes(e.key)) {
        e.preventDefault();
        const subNav = getSubNavigation(navigationState.activeTab, currentProject);
        const subNavItem = subNav[parseInt(e.key) - 1];
        if (subNavItem) {
          updateState({ activeSubTab: subNavItem.id });
        }
      }

      // New project
      if ((e.metaKey || e.ctrlKey) && e.key === 'n') {
        e.preventDefault();
        onNewProject?.();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [navigationState, handleNavigate, onProjectChange, onNewProject, currentProject]);
};

// Responsive breakpoint hook
export const useBreakpoint = () => {
  const [screenSize, setScreenSize] = useState<{
    isMobile: boolean;
    isTablet: boolean;
    isDesktop: boolean;
  }>({
    isMobile: false,
    isTablet: false,
    isDesktop: true
  });

  useEffect(() => {
    const checkScreenSize = () => {
      const width = window.innerWidth;
      setScreenSize({
        isMobile: width < 768,
        isTablet: width >= 768 && width < 1024,
        isDesktop: width >= 1024
      });
    };

    checkScreenSize();
    window.addEventListener('resize', checkScreenSize);
    return () => window.removeEventListener('resize', checkScreenSize);
  }, []);

  return screenSize;
};

// Local storage hook for preferences
export const useNavbarPreferences = () => {
  const [preferences, setPreferences] = useState({
    viewMode: 'table' as 'table' | 'grid',
    favorites: [] as string[],
    sidebarCollapsed: false
  });

  useEffect(() => {
    try {
      const saved = localStorage.getItem('kova-navbar-preferences');
      if (saved) {
        setPreferences(JSON.parse(saved));
      }
    } catch (error) {
      console.warn('Failed to load navbar preferences:', error);
    }
  }, []);

  const updatePreferences = (updates: Partial<typeof preferences>) => {
    const newPreferences = { ...preferences, ...updates };
    setPreferences(newPreferences);

    try {
      localStorage.setItem('kova-navbar-preferences', JSON.stringify(newPreferences));
    } catch (error) {
      console.warn('Failed to save navbar preferences:', error);
    }
  };

  return { preferences, updatePreferences };
};
