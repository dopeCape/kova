import React from 'react';

type IconSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl';
type IconColor = 'primary' | 'secondary' | 'muted' | 'accent' | 'inherit' | 'current';

interface IconProps {
  children: React.ReactElement;
  size?: IconSize;
  color?: IconColor;
  className?: string;
  strokeWidth?: number;
  onClick?: () => void;
  'aria-label'?: string;
  'aria-hidden'?: boolean;
}

// Size mappings
const sizeClasses: Record<IconSize, string> = {
  xs: 'w-3 h-3',
  sm: 'w-4 h-4',
  md: 'w-5 h-5',
  lg: 'w-6 h-6',
  xl: 'w-8 h-8'
};

// Color mappings
const colorClasses: Record<IconColor, string> = {
  primary: 'text-white',
  secondary: 'text-zinc-300',
  muted: 'text-zinc-400',
  accent: 'text-emerald-500',
  inherit: 'text-inherit',
  current: 'text-current'
};

const Icon: React.FC<IconProps> = ({
  children,
  size = 'md',
  color = 'current',
  className = '',
  strokeWidth,
  onClick,
  'aria-label': ariaLabel,
  'aria-hidden': ariaHidden = !ariaLabel,
  ...props
}) => {
  const sizeClass = sizeClasses[size];
  const colorClass = colorClasses[color];

  const baseClasses = 'flex-shrink-0 transition-colors duration-200';
  const interactiveClasses = onClick ? 'cursor-pointer hover:opacity-80' : '';

  const combinedClassName = [
    baseClasses,
    sizeClass,
    colorClass,
    interactiveClasses,
    className
  ].filter(Boolean).join(' ');

  const iconElement = React.cloneElement<any>(children, {
    className: combinedClassName,
    strokeWidth: strokeWidth || (children.props as any).strokeWidth,
    onClick,
    'aria-label': ariaLabel,
    'aria-hidden': ariaHidden,
    role: onClick ? 'button' : undefined,
    tabIndex: onClick ? 0 : undefined,
    onKeyDown: onClick ? (e: React.KeyboardEvent) => {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        onClick();
      }
    } : undefined,
    ...children.props as Object,
    ...props
  });

  return iconElement;
};

// Status indicator icons with predefined colors and animations
interface StatusIconProps extends Omit<IconProps, 'color' | 'children'> {
  status: 'live' | 'building' | 'error' | 'idle';
  animated?: boolean;
}

const StatusIcon: React.FC<StatusIconProps> = ({
  status,
  animated = true,
  size = 'sm',
  className = '',
  ...props
}) => {
  const getStatusIcon = () => {
    const baseClasses = 'rounded-full';
    const animationClass = animated ? {
      live: 'animate-pulse',
      building: 'animate-ping',
      error: '',
      idle: ''
    }[status] : '';

    const statusColor = {
      live: 'bg-emerald-500 shadow-emerald-500/30',
      building: 'bg-zinc-400',
      error: 'bg-zinc-600',
      idle: 'bg-zinc-700'
    }[status];

    const statusSize = {
      xs: 'w-1.5 h-1.5',
      sm: 'w-2 h-2',
      md: 'w-2.5 h-2.5',
      lg: 'w-3 h-3',
      xl: 'w-4 h-4'
    }[size];

    const combinedClassName = [
      baseClasses,
      statusSize,
      statusColor,
      animationClass,
      status === 'live' ? 'shadow-lg' : '',
      className
    ].filter(Boolean).join(' ');

    return <div className={combinedClassName} {...props} />;
  };

  return getStatusIcon();
};


export { StatusIcon, Icon };
export type {
  IconProps,
  IconSize,
  IconColor,
  StatusIconProps,
};
