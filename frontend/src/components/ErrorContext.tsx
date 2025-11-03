import React, { createContext, useContext, useState, ReactNode } from 'react';
import ErrorModal from './ErrorModal';

interface ErrorContextType {
  showError: (message: string, title?: string, details?: string) => void;
}

const ErrorContext = createContext<ErrorContextType | undefined>(undefined);

interface ErrorProviderProps {
  children: ReactNode;
}

export const ErrorProvider: React.FC<ErrorProviderProps> = ({ children }) => {
  const [errorState, setErrorState] = useState<{
    isOpen: boolean;
    title?: string;
    message: string;
    details?: string;
  }>({
    isOpen: false,
    message: '',
  });

  const showError = (message: string, title?: string, details?: string) => {
    setErrorState({
      isOpen: true,
      title: title || 'Error',
      message,
      details,
    });
  };

  const hideError = () => {
    setErrorState(prev => ({ ...prev, isOpen: false }));
  };

  return (
    <ErrorContext.Provider value={{ showError }}>
      {children}
      <ErrorModal
        isOpen={errorState.isOpen}
        onClose={hideError}
        title={errorState.title}
        message={errorState.message}
        details={errorState.details}
      />
    </ErrorContext.Provider>
  );
};

export const useError = () => {
  const context = useContext(ErrorContext);
  if (!context) {
    throw new Error('useError must be used within an ErrorProvider');
  }
  return context;
};