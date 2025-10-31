"use client"
import React, { useRef, useEffect } from 'react';
import { Command, Search, ArrowRight, Globe, Code, Rocket, Terminal, Users } from 'lucide-react';
import { NavigationState, Project } from './types';
import { getContextualCommands } from './navigationConfig';

interface CommandPaletteProps {
  navigationState: NavigationState;
  updateState: (updates: Partial<NavigationState>) => void;
  handleNavigate: (tab: string, subTab?: string) => void;
  currentProject?: string;
  projects: Project[];
  onNewProject: (() => void) | undefined;
}

export const CommandPalette: React.FC<CommandPaletteProps> = ({
  navigationState,
  updateState,
  handleNavigate,
  currentProject,
  projects,
  onNewProject
}) => {
  const commandRef = useRef<HTMLInputElement>(null);

  // Focus input when command palette opens
  useEffect(() => {
    if (navigationState.showCommand && commandRef.current) {
      commandRef.current.focus();
    }
  }, [navigationState.showCommand]);

  if (!navigationState.showCommand) return null;

  const commands = getContextualCommands(
    navigationState.activeTab,
    currentProject,
    handleNavigate,
    onNewProject
  );

  const filteredCommands = commands.filter(cmd =>
    cmd.name.toLowerCase().includes(navigationState.searchQuery.toLowerCase()) ||
    cmd.command.toLowerCase().includes(navigationState.searchQuery.toLowerCase())
  );

  // Search scopes
  const searchScopes = [
    { id: 'global', name: 'Everything', icon: Globe },
    { id: 'projects', name: 'Projects', icon: Code },
    { id: 'deployments', name: 'Deployments', icon: Rocket },
    { id: 'logs', name: 'Logs', icon: Terminal },
    { id: 'team', name: 'Team', icon: Users },
  ];

  // Group commands by section
  const groupedCommands = filteredCommands.reduce((acc, cmd) => {
    const section = cmd.section || 'Other';
    if (!acc[section]) acc[section] = [];
    acc[section].push(cmd);
    return acc;
  }, {} as Record<string, typeof filteredCommands>);

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center pt-[15vh] bg-black/50 backdrop-blur-sm">
      <div className="w-full max-w-2xl mx-4 bg-zinc-900 border border-zinc-700 rounded-lg shadow-2xl">

        {/* Search Header */}
        <div className="flex items-center space-x-3 p-4 border-b border-zinc-700">
          <Command className="w-5 h-5 text-zinc-400" />
          <input
            ref={commandRef}
            type="text"
            placeholder="Type a command or search..."
            value={navigationState.searchQuery}
            onChange={(e) => updateState({ searchQuery: e.target.value })}
            className="flex-1 bg-transparent text-white placeholder-zinc-400 border-0 outline-0 text-lg"
          />

          {/* Search Scope Indicators */}
          <div className="hidden sm:flex items-center space-x-1">
            {searchScopes.map((scope) => (
              <button
                key={scope.id}
                onClick={() => updateState({ searchQuery: scope.name.toLowerCase() + ' ' })}
                className="p-1.5 rounded transition-colors text-zinc-400 hover:text-white"
                title={`Search ${scope.name}`}
              >
                <scope.icon className="w-3 h-3" />
              </button>
            ))}
          </div>

          <kbd className="px-2 py-1 bg-zinc-800 text-zinc-400 rounded text-xs">ESC</kbd>
        </div>

        {/* Command Results */}
        <div className="max-h-80 overflow-y-auto">
          {Object.keys(groupedCommands).length === 0 ? (
            <div className="p-8 text-center text-zinc-400">
              <Search className="w-8 h-8 mx-auto mb-2 opacity-50" />
              <p>No commands found</p>
              <p className="text-xs mt-1">Try searching for actions, pages, or projects</p>
            </div>
          ) : (
            <div className="p-2">
              {Object.entries(groupedCommands).map(([section, commands]) => (
                <div key={section} className="mb-4 last:mb-0">
                  <div className="text-xs text-zinc-400 uppercase tracking-wide px-2 py-2 sticky top-0 bg-zinc-900">
                    {section}
                  </div>
                  {commands.map((cmd, index) => (
                    <button
                      key={index}
                      onClick={() => {
                        cmd.action();
                        updateState({ showCommand: false, searchQuery: '' });
                      }}
                      className="w-full flex items-center justify-between p-3 hover:bg-zinc-800 rounded transition-colors text-left group"
                    >
                      <div>
                        <div className="text-white font-medium">{cmd.name}</div>
                        <div className="text-zinc-400 text-sm">{cmd.command}</div>
                      </div>
                      <div className="flex items-center space-x-2">
                        <kbd className="px-2 py-1 bg-zinc-800 text-zinc-400 rounded text-xs">
                          {cmd.shortcut}
                        </kbd>
                        <ArrowRight className="w-3 h-3 text-zinc-600" />
                      </div>
                    </button>
                  ))}
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Command Tips */}
        <div className="border-t border-zinc-700 p-4 bg-zinc-900/50">
          <div className="flex items-center justify-between text-xs text-zinc-500">
            <div className="flex items-center space-x-4">
              <span><kbd className="bg-zinc-800 px-1 rounded">⌘P</kbd> Projects</span>
              <span><kbd className="bg-zinc-800 px-1 rounded">⌘N</kbd> New</span>
              <span className="hidden sm:inline"><kbd className="bg-zinc-800 px-1 rounded">⌘1-7</kbd> Navigate</span>
            </div>
            <span className="hidden sm:inline">Start typing to search</span>
          </div>
        </div>
      </div>
    </div>
  );
};
