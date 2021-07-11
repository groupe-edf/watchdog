import { ComponentType, FC } from "react";
import { Spinner } from "@chakra-ui/react"
import { Alert } from "./Alert";

interface StatusIndicatorProps {
  error?: Error;
  isLoading?: boolean;
  customErrorMsg?: JSX.Element;
  componentTitle?: string;
}

export const withStatusIndicator = <T extends {}>(Component: ComponentType<T>): FC<StatusIndicatorProps & T> => ({
  error,
  isLoading,
  customErrorMsg,
  componentTitle,
  ...options
}) => {
  if (error) {
    <Alert />
  }
  if (isLoading) {
    return (
      <Spinner
        thickness="4px"
        speed="0.65s"
        emptyColor="gray.200"
        color="brand.100"
        size="xl"/>
    )
  }
  return <Component {...(options as T)} />;
}
