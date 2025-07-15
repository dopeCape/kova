import React from 'react';

type SpinnerVariant = 'spinner' | 'dots' | 'pulse' | 'bars';
type SpinnerSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl';
type SpinnerColor = 'primary' | 'secondary' | 'muted' | 'accent' | 'current';
type SpinnerSpeed = 'slow' | 'normal' | 'fast';

interface SpinnerProps {
  variant?: SpinnerVariant;
  size?: SpinnerSize;
  color?: SpinnerColor;
  speed?: SpinnerSpeed;
  className?: string;
  'aria-label'?: string;
}

interface LoadingProps {
  loading: boolean;
  children: React.ReactNode;
  spinner?: SpinnerProps;
  overlay?: boolean;
  text?: string;
  className?: string;
}

const sizeClasses: Record<SpinnerSize, {
  container: string;
  element: string;
}> = {
  xs: { container: 'w-3 h-3', element: 'w-1 h-1' },
  sm: { container: 'w-4 h-4', element: 'w-1.5 h-1.5' },
  md: { container: 'w-6 h-6', element: 'w-2 h-2' },
  lg: { container: 'w-8 h-8', element: 'w-2.5 h-2.5' },
  xl: { container: 'w-12 h-12', element: 'w-3 h-3' }
};

const colorClasses: Record<SpinnerColor, string> = {
  primary: 'text-white border-white',
  secondary: 'text-zinc-300 border-zinc-300',
  muted: 'text-zinc-400 border-zinc-400',
  accent: 'text-emerald-500 border-emerald-500',
  current: 'text-current border-current'
};

const speedClasses: Record<SpinnerSpeed, string> = {
  slow: '[animation-duration:2s]',
  normal: '[animation-duration:1s]',
  fast: '[animation-duration:0.5s]'
};

const Spinner: React.FC<SpinnerProps> = ({
  variant = 'spinner',
  size = 'md',
  color = 'current',
  speed = 'normal',
  className = '',
  'aria-label': ariaLabel = 'Loading...',
  ...props
}) => {
  const sizeClass = sizeClasses[size];
  const colorClass = colorClasses[color];
  const speedClass = speedClasses[speed];

  const baseClasses = 'inline-block';

  const renderSpinner = () => {
    switch (variant) {
      case 'spinner':
        return (
          <svg
            className={`${sizeClass.container} ${colorClass} animate-spin ${speedClass} ${className}`}
            fill="none"
            viewBox="0 0 24 24"
            role="status"
            aria-label={ariaLabel}
            {...props}
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            />
            <path
              className="opacity-75"
              fill="currentColor"
              d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
        );

      case 'dots':
        return (
          <div
            className={`${baseClasses} ${sizeClass.container} flex items-center justify-center space-x-1 ${className}`}
            role="status"
            aria-label={ariaLabel}
            {...props}
          >
            {[0, 1, 2].map((i) => (
              <div
                key={i}
                className={`${sizeClass.element} ${colorClass.split(' ')[0]} rounded-full animate-pulse ${speedClass}`}
                style={{
                  animationDelay: `${i * 0.2}s`,
                  animationDuration: speed === 'slow' ? '2s' : speed === 'fast' ? '0.5s' : '1s'
                }}
              />
            ))}
          </div>
        );

      case 'pulse':
        return (
          <div
            className={`${baseClasses} ${sizeClass.container} ${colorClass.split(' ')[0]} rounded-full animate-pulse ${speedClass} ${className}`}
            role="status"
            aria-label={ariaLabel}
            {...props}
          />
        );

      case 'bars':
        return (
          <div
            className={`${baseClasses} ${sizeClass.container} flex items-end justify-center space-x-1 ${className}`}
            role="status"
            aria-label={ariaLabel}
            {...props}
          >
            {[0, 1, 2, 3].map((i) => (
              <div
                key={i}
                className={`w-1 ${colorClass.split(' ')[0]} rounded-full animate-pulse`}
                style={{
                  height: `${25 + (i % 2) * 50}%`,
                  animationDelay: `${i * 0.15}s`,
                  animationDuration: speed === 'slow' ? '2s' : speed === 'fast' ? '0.5s' : '1s'
                }}
              />
            ))}
          </div>
        );

      default:
        return null;
    }
  };

  return renderSpinner();
};

// Loading wrapper component
const Loading: React.FC<LoadingProps> = ({
  loading,
  children,
  spinner = {},
  overlay = false,
  text,
  className = '',
  ...props
}) => {
  if (!loading) {
    return <>{children}</>;
  }

  const loadingContent = (
    <div className="flex flex-col items-center justify-center space-y-3">
      <Spinner {...spinner} />
      {text && (
        <div className="text-sm text-zinc-400 font-mono">{text}</div>
      )}
    </div>
  );

  if (overlay) {
    return (
      <div className={`relative ${className}`} {...props}>
        {children}
        <div className="absolute inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-10">
          {loadingContent}
        </div>
      </div>
    );
  }

  return (
    <div className={`flex items-center justify-center p-8 ${className}`} {...props}>
      {loadingContent}
    </div>
  );
};

// Skeleton component for loading placeholders
interface SkeletonProps {
  className?: string;
  width?: string | number;
  height?: string | number;
  rounded?: boolean;
  animate?: boolean;
}

const Skeleton: React.FC<SkeletonProps> = ({
  className = '',
  width,
  height,
  rounded = false,
  animate = true,
  ...props
}) => {
  const baseClasses = 'bg-zinc-800';
  const animationClass = animate ? 'animate-pulse' : '';
  const roundedClass = rounded ? 'rounded-full' : 'rounded';

  const style: React.CSSProperties = {};
  if (width) style.width = typeof width === 'number' ? `${width}px` : width;
  if (height) style.height = typeof height === 'number' ? `${height}px` : height;

  const combinedClassName = [
    baseClasses,
    animationClass,
    roundedClass,
    className
  ].filter(Boolean).join(' ');

  return (
    <div
      className={combinedClassName}
      style={style}
      role="status"
      aria-label="Loading content..."
      {...props}
    />
  );
};

export { Loading, Skeleton, Spinner };
export type {
  SpinnerProps,
  SpinnerVariant,
  SpinnerSize,
  SpinnerColor,
  SpinnerSpeed,
  LoadingProps,
  SkeletonProps
};
