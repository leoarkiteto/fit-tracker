// This file overrides React types to fix compatibility issues in the monorepo
// The mobile app uses React 18.3.x types while PWA uses 18.2.x

import "react";

declare module "react" {
  // Ensure ReactNode includes undefined
  type ReactNode =
    | React.ReactElement
    | string
    | number
    | Iterable<ReactNode>
    | React.ReactPortal
    | boolean
    | null
    | undefined;
}
