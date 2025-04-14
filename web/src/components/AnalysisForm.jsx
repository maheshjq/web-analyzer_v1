import React, { useState } from 'react';

const AnalysisForm = ({ onSubmit, isLoading }) => {
  const [url, setUrl] = useState('');
  const [error, setError] = useState('');

  // Simple URL validation
  const validateUrl = (value) => {
    const urlPattern = /^(https?:\/\/)?(www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)$/;
    return urlPattern.test(value);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    
    // Trim and normalize URL
    const trimmedUrl = url.trim();
    
    // Validate URL
    if (!trimmedUrl) {
      setError('URL is required');
      return;
    }
    
    if (!validateUrl(trimmedUrl)) {
      setError('Please enter a valid URL');
      return;
    }
    
    // Clear error and submit
    setError('');
    onSubmit(trimmedUrl);
  };

  return (
    <div className="mb-8 w-full max-w-2xl mx-auto bg-white rounded-lg shadow-md p-6">
      <h2 className="text-xl font-semibold mb-4 text-gray-800">Enter a URL to analyze</h2>
      
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="url" className="block text-sm font-medium text-gray-700 mb-1">
            Web Page URL
          </label>
          <div className="mt-1 flex rounded-md shadow-sm">
            <input
              type="text"
              id="url"
              name="url"
              placeholder="https://example.com"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              className={`block w-full px-4 py-2 border ${
                error ? 'border-red-500' : 'border-gray-300'
              } rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500`}
              disabled={isLoading}
            />
          </div>
          {error && (
            <p className="mt-2 text-sm text-red-600">{error}</p>
          )}
        </div>
        
        <div className="flex justify-end">
          <button
            type="submit"
            disabled={isLoading}
            className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? "Analyzing..." : "Analyze"}
          </button>
        </div>
      </form>
    </div>
  );
};

export default AnalysisForm;