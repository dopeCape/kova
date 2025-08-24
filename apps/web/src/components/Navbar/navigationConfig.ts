import {
  Activity, Code, Rocket, Settings, BarChart3, Server, Users,
  Eye, Terminal, Monitor, Globe, Shield, FileText, Zap, Clock,
  XCircle, History, Home, Folder
} from 'lucide-react';
import { NavigationItem, BreadcrumbItem, CommandItem } from './types';

export const mainNavItems: NavigationItem[] = [
  { id: 'dashboard', name: 'Dashboard', icon: Activity, shortcut: '1' },
  { id: 'projects', name: 'Projects', icon: Code, shortcut: '2' },
  { id: 'deployments', name: 'Deployments', icon: Rocket, shortcut: '3' },
  { id: 'analytics', name: 'Analytics', icon: BarChart3, shortcut: '4' },
  { id: 'infrastructure', name: 'Infrastructure', icon: Server, shortcut: '5' },
  { id: 'team', name: 'Team', icon: Users, shortcut: '6' },
  { id: 'settings', name: 'Settings', icon: Settings, shortcut: '7' },
];

export const getSubNavigation = (activeTab: string, currentProject?: string): NavigationItem[] => {
  if (activeTab === 'projects' && currentProject) {
    return [
      { id: 'overview', name: 'Overview', icon: Eye, shortcut: '1' },
      { id: 'deployments', name: 'Deployments', icon: Rocket, shortcut: '2', badge: '47' },
      { id: 'analytics', name: 'Analytics', icon: BarChart3, shortcut: '3' },
      { id: 'logs', name: 'Logs', icon: Terminal, shortcut: '4', badge: 'Live', badgeColor: 'emerald' },
      { id: 'settings', name: 'Settings', icon: Settings, shortcut: '5' },
    ];
  }

  if (activeTab === 'deployments') {
    return [
      { id: 'active', name: 'Active', icon: Zap, shortcut: '1', badge: '12', badgeColor: 'emerald' },
      { id: 'history', name: 'History', icon: History, shortcut: '2' },
      { id: 'queued', name: 'Queued', icon: Clock, shortcut: '3', badge: '3', badgeColor: 'yellow' },
      { id: 'failed', name: 'Failed', icon: XCircle, shortcut: '4', badge: '2', badgeColor: 'red' },
    ];
  }

  if (activeTab === 'infrastructure') {
    return [
      { id: 'overview', name: 'Overview', icon: Monitor, shortcut: '1' },
      { id: 'regions', name: 'Regions', icon: Globe, shortcut: '2', badge: '3' },
      { id: 'nodes', name: 'Edge Nodes', icon: Server, shortcut: '3', badge: '47' },
      { id: 'monitoring', name: 'Monitoring', icon: Activity, shortcut: '4' },
    ];
  }

  if (activeTab === 'team') {
    return [
      { id: 'members', name: 'Members', icon: Users, shortcut: '1', badge: '12' },
      { id: 'permissions', name: 'Permissions', icon: Shield, shortcut: '2' },
      { id: 'activity', name: 'Activity', icon: Activity, shortcut: '3' },
      { id: 'billing', name: 'Billing', icon: FileText, shortcut: '4' },
    ];
  }

  return [];
};

export const getBreadcrumbs = (
  activeTab: string,
  activeSubTab: string,
  currentProject?: string
): BreadcrumbItem[] => {
  const base = [{ name: 'Dashboard', path: '/dashboard', icon: Home }];

  if (activeTab === 'projects' && currentProject) {
    return [
      ...base,
      { name: 'Projects', path: '/projects', icon: Code },
      { name: currentProject, path: `/projects/${currentProject}`, icon: Folder },
      {
        name: activeSubTab.charAt(0).toUpperCase() + activeSubTab.slice(1),
        path: `/projects/${currentProject}/${activeSubTab}`,
        icon: Rocket,
        current: true
      }
    ];
  }

  if (activeTab !== 'dashboard') {
    const currentNav = mainNavItems.find(item => item.id === activeTab);
    if (currentNav) {
      return [
        ...base,
        { name: currentNav.name, path: `/${activeTab}`, icon: currentNav.icon, current: true }
      ];
    }
  }

  return [{ ...base[0], current: true }];
};

export const getContextualCommands = (
  activeTab: string,
  currentProject?: string,
  handleNavigate?: (tab: string, subTab?: string) => void,
  onNewProject?: () => void
): CommandItem[] => {
  const base: CommandItem[] = [
    {
      name: 'New Project',
      command: 'new project',
      shortcut: '⌘N',
      action: () => onNewProject?.(),
      section: 'Actions'
    },
    {
      name: 'Go to Dashboard',
      command: 'dashboard',
      shortcut: '⌘1',
      action: () => handleNavigate?.('dashboard'),
      section: 'Navigation'
    },
    {
      name: 'Go to Projects',
      command: 'projects',
      shortcut: '⌘2',
      action: () => handleNavigate?.('projects'),
      section: 'Navigation'
    },
    {
      name: 'Go to Deployments',
      command: 'deployments',
      shortcut: '⌘3',
      action: () => handleNavigate?.('deployments'),
      section: 'Navigation'
    },
    {
      name: 'Go to Analytics',
      command: 'analytics',
      shortcut: '⌘4',
      action: () => handleNavigate?.('analytics'),
      section: 'Navigation'
    },
    {
      name: 'Go to Infrastructure',
      command: 'infrastructure',
      shortcut: '⌘5',
      action: () => handleNavigate?.('infrastructure'),
      section: 'Navigation'
    },
    {
      name: 'Go to Team',
      command: 'team',
      shortcut: '⌘6',
      action: () => handleNavigate?.('team'),
      section: 'Navigation'
    },
  ];

  if (currentProject && activeTab === 'projects') {
    base.unshift(
      {
        name: `Deploy ${currentProject}`,
        command: `deploy ${currentProject}`,
        shortcut: '⌘D',
        action: () => alert('Deploy'),
        section: 'Project Actions'
      },
      {
        name: `${currentProject} Logs`,
        command: `${currentProject} logs`,
        shortcut: '⌘L',
        action: () => handleNavigate?.('projects', 'logs'),
        section: 'Project Actions'
      },
      {
        name: `${currentProject} Analytics`,
        command: `${currentProject} analytics`,
        shortcut: '⌘A',
        action: () => handleNavigate?.('projects', 'analytics'),
        section: 'Project Actions'
      }
    );
  }

  return base;
};
