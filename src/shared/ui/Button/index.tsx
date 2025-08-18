"use client";

import { motion } from 'framer-motion';
import Link from 'next/link';
import { ReactNode } from 'react';

type ButtonProps = {
  children: ReactNode;
  onClick?: () => void;
  href?: string;
  variant?: 'primary' | 'secondary' | 'outline';
  className?: string;
  type?: 'button' | 'submit' | 'reset';
  disabled?: boolean;
  target?: string;
};

export const Button = ({
  children,
  onClick,
  href,
  variant = 'primary',
  className = '',
  type = 'button',
  disabled = false,
  target,
}: ButtonProps) => {
  const baseClasses = 'px-6 py-3 rounded-lg font-medium transition-all duration-300 inline-flex items-center justify-center';
  
  const variantClasses = {
    primary: 'bg-blue-600 text-white hover:bg-blue-700 shadow-lg hover:shadow-xl',
    secondary: 'bg-gray-200 text-gray-800 hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-100 dark:hover:bg-gray-600',
    outline: 'border-2 border-blue-600 text-blue-600 hover:bg-blue-50 dark:hover:bg-gray-800',
  };
  
  const disabledClasses = disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer';
  
  const buttonClasses = `${baseClasses} ${variantClasses[variant]} ${disabledClasses} ${className}`;
  
  const MotionButton = motion.button;
  
  if (href) {
    return (
      <Link 
        href={href} 
        passHref
        className={buttonClasses}
        target={target}
        style={{ display: 'inline-flex' }}
      >
        <motion.span
          className="w-full h-full flex items-center justify-center"
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
        >
          {children}
        </motion.span>
      </Link>
    );
  }
  
  return (
    <MotionButton
      type={type}
      onClick={onClick}
      className={buttonClasses}
      disabled={disabled}
      whileHover={disabled ? {} : { scale: 1.05 }}
      whileTap={disabled ? {} : { scale: 0.95 }}
    >
      {children}
    </MotionButton>
  );
};
