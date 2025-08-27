import { LucideIcon } from 'lucide-react';

export interface User {
  id: string;
  name: string;
  email: string;
  role: string;
  loading: boolean;
}

export interface Project {
  id: string;
  name: string;
  status: 'active' | 'error' | 'building' | 'inactive';
  team: string;
  lastDeploy: string;
  framework?: string;
  repository?: string;
  branch?: string;
  domain?: string;
}

export interface NavigationItem {
  id: string;
  name: string;
  icon: LucideIcon;
  shortcut: string;
  badge?: string;
  badgeColor?: 'emerald' | 'yellow' | 'red' | 'blue' | 'purple';
}

export interface BreadcrumbItem {
  name: string;
  path: string;
  icon: LucideIcon;
  current?: boolean;
}

export interface CommandItem {
  name: string;
  command: string;
  shortcut: string;
  action: () => void;
  section?: string;
}

export interface NavigationState {
  activeTab: string;
  activeSubTab: string;
  showCommand: boolean;
  showProjectSwitcher: boolean;
  showUserMenu: boolean;
  showMobileMenu: boolean;
  searchQuery: string;
  viewMode: 'table' | 'grid';
  favorites: string[];
}

export interface NavbarProps {
  user?: User;
  currentProject?: string;
  projects?: Project[];
  onNavigate?: (tab: string, subTab?: string) => void;
  onProjectChange?: (projectId: string) => void;
  onNewProject?: () => void;
  className?: string;
  notifications?: number;
  systemStatus?: 'operational' | 'degraded' | 'down';
}

export interface SharedNavbarProps {
  navigationState: NavigationState;
  updateState: (updates: Partial<NavigationState>) => void;
  handleNavigate: (tab: string, subTab?: string) => void;
  user: User;
  currentProject?: string;
  projects: Project[];
  onProjectChange?: (projectId: string) => void;
  onNewProject?: () => void;
  notifications: number;
  systemStatus: string;
}

export type ScreenSize = 'mobile' | 'tablet' | 'desktop';
