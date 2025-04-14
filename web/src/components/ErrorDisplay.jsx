import React from 'react';

const ErrorDisplay = ({ error }) => {
  if (!error) return null;

  return (
    <div className="bg-red-50 border-l-4 border-red-500 p-4 my-6 rounded shadow-md">
      <div className="flex">
        <div className="ml-3">
          <h3 className="text-sm font-medium text-red-800">
            Error {error.statusCode && `(${error.statusCode})`}
          </h3>
          <div className="mt-2 text-sm text-red-700">
            <p>{error.message || 'An unknown error occurred'}</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ErrorDisplay;