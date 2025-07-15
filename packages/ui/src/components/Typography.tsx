import React, { JSX } from 'react';

// Types
type TypographyVariant = 'h1' | 'h2' | 'h3' | 'h4' | 'text' | 'code' | 'label';
type TypographySize = 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '3xl';
type TypographyWeight = 'light' | 'normal' | 'medium' | 'semibold' | 'bold';
type TypographyColor = 'primary' | 'secondary' | 'muted' | 'accent' | 'inherit';

interface TypographyProps {
  children: React.ReactNode;
  variant?: TypographyVariant;
  size?: TypographySize;
  weight?: TypographyWeight;
  color?: TypographyColor;
  className?: string;
  as?: keyof JSX.IntrinsicElements;
  mono?: boolean;
}

// Typography variant mappings
const variantElements: Record<TypographyVariant, keyof JSX.IntrinsicElements> = {
  h1: 'h1',
  h2: 'h2',
  h3: 'h3',
  h4: 'h4',
  text: 'p',
  code: 'code',
  label: 'label'
};

const variantSizes: Record<TypographyVariant, TypographySize> = {
  h1: '3xl',
  h2: '2xl',
  h3: 'xl',
  h4: 'lg',
  text: 'md',
  code: 'sm',
  label: 'sm'
};

const variantWeights: Record<TypographyVariant, TypographyWeight> = {
  h1: 'bold',
  h2: 'bold',
  h3: 'semibold',
  h4: 'semibold',
  text: 'normal',
  code: 'normal',
  label: 'medium'
};

// Size mappings
const sizeClasses: Record<TypographySize, string> = {
  xs: 'text-xs',
  sm: 'text-sm',
  md: 'text-base',
  lg: 'text-lg',
  xl: 'text-xl',
  '2xl': 'text-2xl',
  '3xl': 'text-3xl'
};

// Weight mappings
const weightClasses: Record<TypographyWeight, string> = {
  light: 'font-light',
  normal: 'font-normal',
  medium: 'font-medium',
  semibold: 'font-semibold',
  bold: 'font-bold'
};

// Color mappings
const colorClasses: Record<TypographyColor, string> = {
  primary: 'text-white',
  secondary: 'text-zinc-300',
  muted: 'text-zinc-400',
  accent: 'text-emerald-500',
  inherit: 'text-inherit'
};

// Special styles for variants
const variantStyles: Record<TypographyVariant, string> = {
  h1: 'tracking-tight leading-tight',
  h2: 'tracking-tight leading-tight',
  h3: 'tracking-tight leading-snug',
  h4: 'tracking-tight leading-snug',
  text: 'leading-relaxed',
  code: 'font-mono bg-zinc-900/50 px-1.5 py-0.5 rounded text-zinc-300 border border-zinc-800',
  label: 'leading-none'
};

const Typography: React.FC<TypographyProps> = ({
  children,
  variant = 'text',
  size,
  weight,
  color = 'primary',
  className = '',
  as,
  mono = false,
  ...props
}) => {
  // Determine the HTML element to render
  const Element = as || variantElements[variant];

  // Determine default size and weight based on variant
  const finalSize = size || variantSizes[variant];
  const finalWeight = weight || variantWeights[variant];

  // Build className
  const baseClasses = 'transition-colors duration-200';
  const sizeClass = sizeClasses[finalSize];
  const weightClass = weightClasses[finalWeight];
  const colorClass = colorClasses[color];
  const variantClass = variantStyles[variant];
  const monoClass = mono || variant === 'code' ? 'font-mono' : '';

  const combinedClassName = [
    baseClasses,
    sizeClass,
    weightClass,
    colorClass,
    variantClass,
    monoClass,
    className
  ].filter(Boolean).join(' ');

  return (
    <Element className={combinedClassName} {...props}>
      {children}
    </Element>
  );
};

// Convenience components for common use cases
const Heading: React.FC<Omit<TypographyProps, 'variant'> & { level?: 1 | 2 | 3 | 4 }> = ({
  level = 1,
  ...props
}) => (
  <Typography variant={`h${level}` as TypographyVariant} {...props} />
);

const Text: React.FC<Omit<TypographyProps, 'variant'>> = (props) => (
  <Typography variant="text" {...props} />
);

const Code: React.FC<Omit<TypographyProps, 'variant'>> = (props) => (
  <Typography variant="code" {...props} />
);

const Label: React.FC<Omit<TypographyProps, 'variant'>> = (props) => (
  <Typography variant="label" {...props} />
);

export { Typography };
export { Heading, Text, Code, Label };
export type {
  TypographyProps,
  TypographyVariant,
  TypographySize,
  TypographyWeight,
  TypographyColor
};
