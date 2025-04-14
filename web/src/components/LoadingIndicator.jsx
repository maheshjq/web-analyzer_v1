import React from 'react';

const LoadingIndicator = () => {
  return (
    <div className="flex justify-center items-center my-12">
      <div className="flex flex-col items-center">
        <div className="w-16 h-16 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
        <span className="mt-4 text-lg text-gray-700">Analyzing web page...</span>
      </div>
    </div>
  );
};

export default LoadingIndicator;