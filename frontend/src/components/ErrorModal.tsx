import { AlertTriangle, X } from 'lucide-react';
import React from 'react';

interface ErrorModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  message: string;
  details?: string;
}

const ErrorModal: React.FC<ErrorModalProps> = ({ 
  isOpen, 
  onClose, 
  title = 'Error', 
  message, 
  details 
}) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div 
        className="absolute inset-0 bg-black/70 backdrop-blur-sm" 
        onClick={onClose}
      />
      
      {/* Modal Content */}
      <div className="relative bg-neutral-800 border border-red-600 rounded-lg shadow-2xl max-w-md w-full max-h-[90vh] overflow-y-auto z-10">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-neutral-700">
          <div className="flex items-center gap-2">
            <AlertTriangle className="h-5 w-5 text-red-500" />
            <h3 className="text-lg font-semibold text-red-400">{title}</h3>
          </div>
          
          <button
            onClick={onClose}
            className="text-neutral-400 hover:text-white transition-colors"
            aria-label="Close"
          >
            <X className="h-5 w-5" />
          </button>
        </div>
        
        {/* Body */}
        <div className="p-4">
          <p className="text-neutral-200 mb-3">{message}</p>
          
          {details && (
            <div className="mt-3">
              <h4 className="text-sm font-medium text-neutral-300 mb-1">Details:</h4>
              <pre className="text-sm text-neutral-400 bg-neutral-900 p-2 rounded border border-neutral-700 overflow-x-auto">
                {details}
              </pre>
            </div>
          )}
        </div>
        
        {/* Footer */}
        <div className="p-4 border-t border-neutral-700 flex justify-end">
          <button
            onClick={onClose}
            className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-md transition-colors focus:outline-none focus:ring-2 focus:ring-red-500"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
};

export default ErrorModal;