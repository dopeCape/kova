import React from 'react';

// Types
type DividerVariant = 'solid' | 'dashed' | 'dotted' | 'gradient';
type DividerOrientation = 'horizontal' | 'vertical';
type DividerSize = 'sm' | 'md' | 'lg';
type DividerColor = 'default' | 'muted' | 'accent';

interface DividerProps {
  variant?: DividerVariant;
  orientation?: DividerOrientation;
  size?: DividerSize;
  color?: DividerColor;
  className?: string;
  children?: React.ReactNode;
  spacing?: 'none' | 'sm' | 'md' | 'lg' | 'xl';
}

// Color mappings
const colorClasses: Record<DividerColor, string> = {
  default: 'border-zinc-800',
  muted: 'border-zinc-900',
  accent: 'border-emerald-500/30'
};

// Size mappings for thickness
const sizeClasses: Record<DividerOrientation, Record<DividerSize, string>> = {
  horizontal: {
    sm: 'border-t',
    md: 'border-t-2',
    lg: 'border-t-4'
  },
  vertical: {
    sm: 'border-l',
    md: 'border-l-2',
    lg: 'border-l-4'
  }
};

// Spacing mappings
const spacingClasses: Record<DividerOrientation, Record<string, string>> = {
  horizontal: {
    none: '',
    sm: 'my-2',
    md: 'my-4',
    lg: 'my-6',
    xl: 'my-8'
  },
  vertical: {
    none: '',
    sm: 'mx-2',
    md: 'mx-4',
    lg: 'mx-6',
    xl: 'mx-8'
  }
};

// Variant styles
const variantClasses: Record<DividerVariant, string> = {
  solid: 'border-solid',
  dashed: 'border-dashed',
  dotted: 'border-dotted',
  gradient: ''
};

const Divider: React.FC<DividerProps> = ({
  variant = 'solid',
  orientation = 'horizontal',
  size = 'sm',
  color = 'default',
  className = '',
  children,
  spacing = 'md',
  ...props
}) => {
  const colorClass = colorClasses[color];
  const sizeClass = sizeClasses[orientation][size];
  const variantClass = variantClasses[variant];
  const spacingClass = spacingClasses[orientation][spacing];

  // Base classes
  const baseClasses = orientation === 'horizontal'
    ? 'w-full'
    : 'h-full min-h-4';

  // Handle gradient variant separately
  if (variant === 'gradient') {
    const gradientClass = orientation === 'horizontal'
      ? 'h-px bg-gradient-to-r from-transparent via-zinc-800 to-transparent'
      : 'w-px bg-gradient-to-b from-transparent via-zinc-800 to-transparent';

    const combinedClassName = [
      baseClasses,
      gradientClass,
      spacingClass,
      className
    ].filter(Boolean).join(' ');

    return (
      <div className={combinedClassName} role="separator" {...props}>
        {children}
      </div>
    );
  }

  // Handle divider with text/content
  if (children) {
    if (orientation === 'vertical') {
      // Vertical divider with content is complex, just render simple divider
      const combinedClassName = [
        baseClasses,
        sizeClass,
        variantClass,
        colorClass,
        spacingClass,
        className
      ].filter(Boolean).join(' ');

      return (
        <div className={combinedClassName} role="separator" {...props} />
      );
    }

    // Horizontal divider with text
    return (
      <div className={`relative flex items-center ${spacingClass} ${className}`} {...props}>
        <div className={`flex-1 ${sizeClass} ${variantClass} ${colorClass}`} />
        <div className="relative flex justify-center text-xs">
          <span className="bg-black px-4 text-zinc-500 font-mono uppercase tracking-wider">
            {children}
          </span>
        </div>
        <div className={`flex-1 ${sizeClass} ${variantClass} ${colorClass}`} />
      </div>
    );
  }

  // Simple divider without content
  const combinedClassName = [
    baseClasses,
    sizeClass,
    variantClass,
    colorClass,
    spacingClass,
    className
  ].filter(Boolean).join(' ');

  return (
    <div className={combinedClassName} role="separator" {...props} />
  );
};

// Convenience components for common patterns
const HDivider: React.FC<Omit<DividerProps, 'orientation'>> = (props) => (
  <Divider orientation="horizontal" {...props} />
);

const VDivider: React.FC<Omit<DividerProps, 'orientation'>> = (props) => (
  <Divider orientation="vertical" {...props} />
);

// Section divider with semantic meaning
interface SectionDividerProps extends Omit<DividerProps, 'children'> {
  title?: string;
  subtitle?: string;
}

const SectionDivider: React.FC<SectionDividerProps> = ({
  title,
  subtitle,
  ...props
}) => {
  if (!title && !subtitle) {
    return <Divider {...props} />;
  }

  return (
    <div className="relative">
      <Divider {...props} />
      {(title || subtitle) && (
        <div className="absolute left-0 top-1/2 transform -translate-y-1/2 bg-black pr-4">
          {title && (
            <h3 className="text-sm font-medium text-zinc-300 mb-1">{title}</h3>
          )}
          {subtitle && (
            <p className="text-xs text-zinc-500 font-mono">{subtitle}</p>
          )}
        </div>
      )}
    </div>
  );
};

// Dot separator for inline content
interface DotSeparatorProps {
  className?: string;
  color?: DividerColor;
}

const DotSeparator: React.FC<DotSeparatorProps> = ({
  className = '',
  color = 'default',
  ...props
}) => {
  const colorMap = {
    default: 'bg-zinc-600',
    muted: 'bg-zinc-700',
    accent: 'bg-emerald-500'
  };

  return (
    <div
      className={`w-1 h-1 rounded-full ${colorMap[color]} ${className}`}
      role="separator"
      {...props}
    />
  );
};

export { HDivider, VDivider, SectionDivider, DotSeparator, Divider };
export type {
  DividerProps,
  DividerVariant,
  DividerOrientation,
  DividerSize,
  DividerColor,
  SectionDividerProps,
  DotSeparatorProps
};
