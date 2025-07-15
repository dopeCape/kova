import React from 'react';
import { Loader2 } from 'lucide-react';

// Types
type ButtonVariant = 'primary' | 'secondary' | 'ghost';
type ButtonSize = 'sm' | 'md' | 'lg';
type IconPosition = 'left' | 'right';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children?: React.ReactNode;
  variant?: ButtonVariant;
  size?: ButtonSize;
  loading?: boolean;
  disabled?: boolean;
  className?: string;
  icon?: React.ReactElement;
  iconPosition?: IconPosition;
  fullWidth?: boolean;
}

// Button variant styles
const buttonVariants: Record<ButtonVariant, {
  base: string;
  hover: string;
  active: string;
  disabled: string;
}> = {
  primary: {
    base: "bg-emerald-500 text-white border-emerald-500 shadow-lg shadow-emerald-500/20",
    hover: "hover:bg-emerald-400 hover:border-emerald-400 hover:shadow-emerald-500/30",
    active: "active:bg-emerald-600 active:border-emerald-600",
    disabled: "disabled:bg-emerald-500/50 disabled:border-emerald-500/50 disabled:shadow-none disabled:cursor-not-allowed"
  },
  secondary: {
    base: "bg-transparent text-zinc-300 border-zinc-700",
    hover: "hover:bg-zinc-800/50 hover:text-white hover:border-zinc-600",
    active: "active:bg-zinc-800 active:border-zinc-500",
    disabled: "disabled:text-zinc-500 disabled:border-zinc-800 disabled:cursor-not-allowed"
  },
  ghost: {
    base: "bg-transparent text-zinc-400 border-transparent",
    hover: "hover:bg-zinc-800/30 hover:text-zinc-200",
    active: "active:bg-zinc-800/50 active:text-white",
    disabled: "disabled:text-zinc-600 disabled:cursor-not-allowed"
  }
};

// Button sizes
const buttonSizes: Record<ButtonSize, string> = {
  sm: "px-3 py-2 text-xs",
  md: "px-4 py-3 text-sm",
  lg: "px-6 py-4 text-base"
};

const Button: React.FC<ButtonProps> = ({
  children,
  variant = 'primary',
  size = 'md',
  loading = false,
  disabled = false,
  className = '',
  icon,
  iconPosition = 'left',
  fullWidth = false,
  type = 'button',
  ...props
}) => {
  const baseStyles = "relative inline-flex items-center justify-center font-medium border rounded-xl transition-all duration-200 transform focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2 focus:ring-offset-black";

  const variantStyles = buttonVariants[variant];
  const sizeStyles = buttonSizes[size];

  const scaleEffect = variant === 'primary' ? "hover:scale-105" : "hover:scale-[1.02]";

  const combinedClassName = [
    baseStyles,
    variantStyles.base,
    variantStyles.hover,
    variantStyles.active,
    variantStyles.disabled,
    sizeStyles,
    scaleEffect,
    fullWidth ? 'w-full' : '',
    className
  ].filter(Boolean).join(' ');

  const isDisabled = disabled || loading;

  const renderIcon = (): React.ReactElement | null => {
    if (loading) {
      return <Loader2 className="w-4 h-4 animate-spin" />;
    }
    if (icon) {
      return React.cloneElement<any>(icon, {
        className: "w-4 h-4",
        ...icon.props as Object
      });
    }
    return null;
  };

  const hasIcon = loading || !!icon;
  const iconSpacing = size === 'sm' ? 'space-x-1.5' : size === 'lg' ? 'space-x-3' : 'space-x-2';

  return (
    <button
      type={type}
      disabled={isDisabled}
      className={combinedClassName}
      {...props}
    >
      {/* Subtle shimmer effect for primary buttons */}
      {variant === 'primary' && !isDisabled && (
        <div className="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity rounded-xl pointer-events-none" />
      )}

      {/* Content wrapper */}
      <span className={`relative flex items-center ${hasIcon ? iconSpacing : ''}`}>
        {hasIcon && iconPosition === 'left' && renderIcon()}
        {children && <span>{children}</span>}
        {hasIcon && iconPosition === 'right' && renderIcon()}
      </span>
    </button>
  );
};

export type { ButtonProps, ButtonVariant, ButtonSize, IconPosition, };
export { Button };
