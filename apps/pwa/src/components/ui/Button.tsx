"use client";

import React from "react";
import { Loader2 } from "lucide-react";

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "danger" | "ghost" | "outline";
  size?: "sm" | "md" | "lg";
  loading?: boolean;
  icon?: React.ReactNode;
}

export function Button({
  children,
  variant = "primary",
  size = "md",
  loading = false,
  icon,
  className = "",
  disabled,
  ...props
}: ButtonProps) {
  const baseStyles =
    "inline-flex items-center justify-center gap-2 font-semibold rounded-xl transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed";

  const variants = {
    primary:
      "bg-primary-500 hover:bg-primary-600 text-white shadow-md hover:shadow-lg hover:-translate-y-0.5 active:translate-y-0 focus:ring-primary-500",
    secondary:
      "bg-surface-200 hover:bg-surface-300 text-surface-900 focus:ring-surface-400",
    danger:
      "bg-error hover:bg-red-600 text-white focus:ring-red-500",
    ghost:
      "bg-transparent hover:bg-surface-200 text-surface-600 hover:text-surface-900",
    outline:
      "bg-transparent border-2 border-primary-500 text-primary-500 hover:bg-primary-50 focus:ring-primary-500",
  };

  const sizes = {
    sm: "py-2 px-4 text-sm",
    md: "py-3 px-5 text-base",
    lg: "py-3.5 px-6 text-lg",
  };

  return (
    <button
      className={`${baseStyles} ${variants[variant]} ${sizes[size]} ${className}`}
      disabled={disabled || loading}
      {...props}
    >
      {loading ? (
        <Loader2 className="w-5 h-5 animate-spin" />
      ) : icon ? (
        icon
      ) : null}
      {children}
    </button>
  );
}
