import React, { useState, forwardRef } from 'react';
import { Eye, EyeOff, Search, X } from 'lucide-react';

type InputVariant = 'default' | 'search' | 'password';
type InputSize = 'sm' | 'md' | 'lg';
type InputState = 'default' | 'error' | 'success';

interface BaseInputProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'> {
  variant?: InputVariant;
  size?: InputSize;
  state?: InputState;
  label?: string;
  hint?: string;
  error?: string;
  success?: string;
  className?: string;
  containerClassName?: string;
  leftIcon?: React.ReactElement;
  rightIcon?: React.ReactElement;
  clearable?: boolean;
  onClear?: () => void;
}

interface InputProps extends BaseInputProps {
  variant?: 'default';
}

interface SearchInputProps extends BaseInputProps {
  variant: 'search';
  onSearch?: (value: string) => void;
  searchIcon?: boolean;
}

interface PasswordInputProps extends BaseInputProps {
  variant: 'password';
  showPasswordToggle?: boolean;
}

type AllInputProps = InputProps | SearchInputProps | PasswordInputProps;

// Size mappings
const sizeClasses: Record<InputSize, {
  input: string;
  icon: string;
  text: string;
}> = {
  sm: {
    input: 'px-3 py-2 text-sm',
    icon: 'w-4 h-4',
    text: 'text-xs'
  },
  md: {
    input: 'px-4 py-3 text-sm',
    icon: 'w-4 h-4',
    text: 'text-sm'
  },
  lg: {
    input: 'px-4 py-4 text-base',
    icon: 'w-5 h-5',
    text: 'text-sm'
  }
};

// State mappings
const stateClasses: Record<InputState, {
  border: string;
  focus: string;
  text: string;
}> = {
  default: {
    border: 'border-zinc-800',
    focus: 'focus:border-emerald-500',
    text: 'text-zinc-400'
  },
  error: {
    border: 'border-red-500/50',
    focus: 'focus:border-red-500',
    text: 'text-red-400'
  },
  success: {
    border: 'border-emerald-500/50',
    focus: 'focus:border-emerald-500',
    text: 'text-emerald-400'
  }
};

const Input = forwardRef<HTMLInputElement, AllInputProps>(({
  variant = 'default',
  size = 'md',
  state = 'default',
  label,
  hint,
  error,
  success,
  className = '',
  containerClassName = '',
  leftIcon,
  rightIcon,
  clearable = false,
  onClear,
  ...props
}, ref) => {
  const [showPassword, setShowPassword] = useState(false);
  const [internalValue, setInternalValue] = useState(props.value || '');

  const sizeClass = sizeClasses[size];
  const stateClass = stateClasses[state];

  // Determine final state based on error/success props
  const finalState = error ? 'error' : success ? 'success' : state;
  const finalStateClass = stateClasses[finalState];

  // Base input classes
  const baseInputClasses = [
    'w-full bg-zinc-900/80 backdrop-blur-sm border rounded-xl',
    'text-white placeholder-zinc-500',
    'focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2 focus:ring-offset-black',
    'transition-all duration-200',
    'font-mono',
    sizeClass.input,
    finalStateClass.border,
    finalStateClass.focus
  ].join(' ');

  // Handle password visibility
  const inputType = variant === 'password'
    ? (showPassword ? 'text' : 'password')
    : props.type || 'text';

  // Handle icons and special functionality
  const renderLeftIcon = () => {
    if (variant === 'search' && (props as SearchInputProps).searchIcon !== false) {
      return <Search className={`${sizeClass.icon} text-zinc-500`} />;
    }
    if (leftIcon) {
      return React.cloneElement<any>(leftIcon, {
        className: `${sizeClass.icon} text-zinc-500 ${(leftIcon.props as any).className || ''}`
      });
    }
    return null;
  };

  const renderRightIcon = () => {
    const icons = [];

    // Clear button
    if (clearable && internalValue) {
      icons.push(
        <button
          key="clear"
          type="button"
          onClick={() => {
            setInternalValue('');
            onClear?.();
            if (props.onChange) {
              props.onChange({ target: { value: '' } } as React.ChangeEvent<HTMLInputElement>);
            }
          }}
          className="text-zinc-500 hover:text-zinc-300 transition-colors"
        >
          <X className={sizeClass.icon} />
        </button>
      );
    }

    // Password toggle
    if (variant === 'password' && (props as PasswordInputProps).showPasswordToggle !== false) {
      icons.push(
        <button
          key="password-toggle"
          type="button"
          onClick={() => setShowPassword(!showPassword)}
          className="text-zinc-500 hover:text-emerald-500 transition-colors"
        >
          {showPassword ? <EyeOff className={sizeClass.icon} /> : <Eye className={sizeClass.icon} />}
        </button>
      );
    }

    // Custom right icon
    if (rightIcon) {
      icons.push(
        React.cloneElement<any>(rightIcon, {
          key: 'custom-right',
          className: `${sizeClass.icon} text-zinc-500 ${(rightIcon.props as any).className || ''}`
        })
      );
    }

    return icons.length > 0 ? (
      <div className="flex items-center space-x-2">
        {icons}
      </div>
    ) : null;
  };

  const hasLeftIcon = renderLeftIcon() !== null;
  const hasRightIcon = renderRightIcon() !== null;

  // Padding adjustments for icons
  const paddingClass = [
    hasLeftIcon ? (size === 'lg' ? 'pl-12' : 'pl-10') : '',
    hasRightIcon ? (size === 'lg' ? 'pr-12' : 'pr-10') : ''
  ].filter(Boolean).join(' ');

  const inputClassName = [
    baseInputClasses,
    paddingClass,
    className
  ].filter(Boolean).join(' ');

  // Handle input changes
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInternalValue(e.target.value);

    // Handle search functionality
    if (variant === 'search' && (props as SearchInputProps).onSearch) {
      (props as SearchInputProps).onSearch?.(e.target.value);
    }

    props.onChange?.(e);
  };

  return (
    <div className={`space-y-2 ${containerClassName}`}>
      {/* Label */}
      {label && (
        <label className={`block font-medium text-zinc-300 ${sizeClass.text}`}>
          {label}
        </label>
      )}

      {/* Input container */}
      <div className="relative">
        {/* Left icon */}
        {hasLeftIcon && (
          <div className={`absolute left-3 top-1/2 transform -translate-y-1/2 pointer-events-none ${size === 'lg' ? 'left-4' : ''}`}>
            {renderLeftIcon()}
          </div>
        )}

        {/* Input */}
        <input
          ref={ref}
          type={inputType}
          className={inputClassName}
          value={internalValue}
          onChange={handleChange}
          {...props}
        />

        {/* Right icons */}
        {hasRightIcon && (
          <div className={`absolute right-3 top-1/2 transform -translate-y-1/2 ${size === 'lg' ? 'right-4' : ''}`}>
            {renderRightIcon()}
          </div>
        )}

        {/* Focus overlay effect */}
        <div className="absolute inset-0 rounded-xl bg-emerald-500/5 opacity-0 focus-within:opacity-100 transition-opacity pointer-events-none" />
      </div>

      {/* Helper text */}
      {(hint || error || success) && (
        <div className={`${sizeClass.text} ${finalStateClass.text} font-mono`}>
          {error || success || hint}
        </div>
      )}
    </div>
  );
});

Input.displayName = 'Input';

export { Input };
export type {
  AllInputProps as InputProps,
  BaseInputProps,
  SearchInputProps,
  PasswordInputProps,
  InputVariant,
  InputSize,
  InputState
};
