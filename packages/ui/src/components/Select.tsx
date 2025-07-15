import React, { useState, useRef, useEffect, forwardRef } from 'react';
import { ChevronDown, Check, X, Search } from 'lucide-react';

// Types
type SelectSize = 'sm' | 'md' | 'lg';
type SelectState = 'default' | 'error' | 'success';

interface SelectOption {
  value: string | number;
  label: string;
  disabled?: boolean;
  icon?: React.ReactElement;
  description?: string;
}

interface SelectProps {
  options: SelectOption[];
  value?: string | number;
  defaultValue?: string | number;
  placeholder?: string;
  size?: SelectSize;
  state?: SelectState;
  label?: string;
  hint?: string;
  error?: string;
  success?: string;
  disabled?: boolean;
  multiple?: boolean;
  searchable?: boolean;
  clearable?: boolean;
  className?: string;
  containerClassName?: string;
  maxHeight?: string;
  onChange?: (value: string | number | (string | number)[]) => void;
  onSearch?: (searchTerm: string) => void;
}

// Size mappings
const sizeClasses: Record<SelectSize, {
  trigger: string;
  text: string;
  icon: string;
}> = {
  sm: {
    trigger: 'px-3 py-2 text-sm',
    text: 'text-xs',
    icon: 'w-4 h-4'
  },
  md: {
    trigger: 'px-4 py-3 text-sm',
    text: 'text-sm',
    icon: 'w-4 h-4'
  },
  lg: {
    trigger: 'px-4 py-4 text-base',
    text: 'text-sm',
    icon: 'w-5 h-5'
  }
};

// State mappings
const stateClasses: Record<SelectState, {
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

const Select = forwardRef<HTMLDivElement, SelectProps>(({
  options,
  value,
  defaultValue,
  placeholder = 'Select an option...',
  size = 'md',
  state = 'default',
  label,
  hint,
  error,
  success,
  disabled = false,
  multiple = false,
  searchable = false,
  clearable = false,
  className = '',
  containerClassName = '',
  maxHeight = '16rem',
  onChange,
  onSearch,
  ...props
}, ref) => {
  const [isOpen, setIsOpen] = useState(false);
  const [selectedValues, setSelectedValues] = useState<(string | number)[]>(
    multiple
      ? (Array.isArray(value) ? value : defaultValue ? [defaultValue] : [])
      : (value !== undefined ? [value] : defaultValue ? [defaultValue] : [])
  );
  const [searchTerm, setSearchTerm] = useState('');
  const [focusedIndex, setFocusedIndex] = useState(-1);

  const triggerRef = useRef<HTMLButtonElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);

  const sizeClass = sizeClasses[size];
  const finalState = error ? 'error' : success ? 'success' : state;
  const stateClass = stateClasses[finalState];

  // Filter options based on search
  const filteredOptions = searchable && searchTerm
    ? options.filter(option =>
      option.label.toLowerCase().includes(searchTerm.toLowerCase())
    )
    : options;

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node) &&
        !triggerRef.current?.contains(event.target as Node)
      ) {
        setIsOpen(false);
        setSearchTerm('');
        setFocusedIndex(-1);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Focus search input when dropdown opens
  useEffect(() => {
    if (isOpen && searchable && searchInputRef.current) {
      searchInputRef.current.focus();
    }
  }, [isOpen, searchable]);

  // Handle keyboard navigation
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!isOpen) {
      if (e.key === 'Enter' || e.key === ' ' || e.key === 'ArrowDown') {
        e.preventDefault();
        setIsOpen(true);
        setFocusedIndex(0);
      }
      return;
    }

    switch (e.key) {
      case 'Escape':
        e.preventDefault();
        setIsOpen(false);
        setFocusedIndex(-1);
        triggerRef.current?.focus();
        break;
      case 'ArrowDown':
        e.preventDefault();
        setFocusedIndex(prev =>
          prev < filteredOptions.length - 1 ? prev + 1 : 0
        );
        break;
      case 'ArrowUp':
        e.preventDefault();
        setFocusedIndex(prev =>
          prev > 0 ? prev - 1 : filteredOptions.length - 1
        );
        break;
      case 'Enter':
        e.preventDefault();
        if (focusedIndex >= 0 && focusedIndex < filteredOptions.length) {
          handleOptionSelect(filteredOptions[focusedIndex]);
        }
        break;
    }
  };

  const handleOptionSelect = (option: SelectOption) => {
    if (option.disabled) return;

    let newSelectedValues: (string | number)[];

    if (multiple) {
      if (selectedValues.includes(option.value)) {
        newSelectedValues = selectedValues.filter(v => v !== option.value);
      } else {
        newSelectedValues = [...selectedValues, option.value];
      }
    } else {
      newSelectedValues = [option.value];
      setIsOpen(false);
    }

    setSelectedValues(newSelectedValues);
    onChange?.(multiple ? newSelectedValues : newSelectedValues[0]);

    if (!multiple) {
      setSearchTerm('');
      setFocusedIndex(-1);
    }
  };

  const handleClear = (e: React.MouseEvent) => {
    e.stopPropagation();
    setSelectedValues([]);
    onChange?.(multiple ? [] : '');
  };

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setSearchTerm(value);
    onSearch?.(value);
    setFocusedIndex(-1);
  };

  const getSelectedLabel = () => {
    if (selectedValues.length === 0) return placeholder;

    if (multiple) {
      if (selectedValues.length === 1) {
        const option = options.find(opt => opt.value === selectedValues[0]);
        return option?.label || '';
      }
      return `${selectedValues.length} selected`;
    }

    const option = options.find(opt => opt.value === selectedValues[0]);
    return option?.label || '';
  };

  const triggerClasses = [
    'relative w-full bg-zinc-900/80 backdrop-blur-sm border rounded-xl',
    'text-left text-white cursor-pointer',
    'focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2 focus:ring-offset-black',
    'transition-all duration-200',
    'flex items-center justify-between',
    sizeClass.trigger,
    stateClass.border,
    stateClass.focus,
    disabled ? 'opacity-50 cursor-not-allowed' : 'hover:border-zinc-700',
    className
  ].filter(Boolean).join(' ');

  return (
    <div className={`relative ${containerClassName}`} ref={ref}>
      {/* Label */}
      {label && (
        <label className={`block font-medium text-zinc-300 mb-2 ${sizeClass.text}`}>
          {label}
        </label>
      )}

      {/* Trigger */}
      <button
        ref={triggerRef}
        type="button"
        className={triggerClasses}
        onClick={() => !disabled && setIsOpen(!isOpen)}
        onKeyDown={handleKeyDown}
        disabled={disabled}
        aria-haspopup="listbox"
        aria-expanded={isOpen}
        {...props}
      >
        <span className={selectedValues.length === 0 ? 'text-zinc-500' : ''}>
          {getSelectedLabel()}
        </span>

        <div className="flex items-center space-x-2">
          {clearable && selectedValues.length > 0 && !disabled && (
            <button
              onClick={handleClear}
              className="text-zinc-500 hover:text-zinc-300 transition-colors"
            >
              <X className={sizeClass.icon} />
            </button>
          )}
          <ChevronDown
            className={`${sizeClass.icon} text-zinc-500 transition-transform duration-200 ${isOpen ? 'rotate-180' : ''
              }`}
          />
        </div>

        {/* Focus overlay */}
        <div className="absolute inset-0 rounded-xl bg-emerald-500/5 opacity-0 focus-within:opacity-100 transition-opacity pointer-events-none" />
      </button>

      {/* Dropdown */}
      {isOpen && (
        <div
          ref={dropdownRef}
          className="absolute z-50 w-full mt-2 bg-zinc-900 border border-zinc-800 rounded-xl shadow-2xl backdrop-blur-xl"
          style={{ maxHeight }}
        >
          {/* Search input */}
          {searchable && (
            <div className="p-3 border-b border-zinc-800">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-zinc-500" />
                <input
                  ref={searchInputRef}
                  type="text"
                  value={searchTerm}
                  onChange={handleSearchChange}
                  placeholder="Search options..."
                  className="w-full pl-10 pr-4 py-2 bg-zinc-800 border border-zinc-700 rounded-lg text-white placeholder-zinc-500 focus:outline-none focus:border-emerald-500 transition-colors text-sm"
                />
              </div>
            </div>
          )}

          {/* Options */}
          <div className="max-h-48 overflow-y-auto py-2">
            {filteredOptions.length === 0 ? (
              <div className="px-4 py-3 text-sm text-zinc-500 text-center font-mono">
                No options found
              </div>
            ) : (
              filteredOptions.map((option, index) => {
                const isSelected = selectedValues.includes(option.value);
                const isFocused = index === focusedIndex;

                return (
                  <div
                    key={option.value}
                    className={`
                      px-4 py-3 cursor-pointer transition-colors flex items-center justify-between
                      ${isFocused ? 'bg-zinc-800' : 'hover:bg-zinc-800/50'}
                      ${option.disabled ? 'opacity-50 cursor-not-allowed' : ''}
                    `}
                    onClick={() => handleOptionSelect(option)}
                  >
                    <div className="flex items-center space-x-3 flex-1">
                      {option.icon && (
                        <span className="text-zinc-400">
                          {React.cloneElement<any>(option.icon, { className: sizeClass.icon })}
                        </span>
                      )}
                      <div className="flex-1">
                        <div className="text-white text-sm">{option.label}</div>
                        {option.description && (
                          <div className="text-zinc-500 text-xs font-mono mt-1">
                            {option.description}
                          </div>
                        )}
                      </div>
                    </div>

                    {isSelected && (
                      <Check className={`${sizeClass.icon} text-emerald-500`} />
                    )}
                  </div>
                );
              })
            )}
          </div>
        </div>
      )}

      {/* Helper text */}
      {(hint || error || success) && (
        <div className={`mt-2 ${sizeClass.text} ${stateClass.text} font-mono`}>
          {error || success || hint}
        </div>
      )}
    </div>
  );
});

Select.displayName = 'Select';

export { Select };
export type { SelectProps, SelectOption, SelectSize, SelectState };

