import React, { useState } from 'react';
import axios from 'axios';

// API service function - moved inline to avoid dependency issues
const analyzeWebPage = async (url) => {
  try {
    const API_URL = process.env.REACT_APP_API_URL || '/api';
    const response = await axios.post(`${API_URL}/analyze`, { url });
    return response.data;
  } catch (error) {
    if (error.response) {
      // Server responded with an error
      throw error.response.data;
    } else if (error.request) {
      // Request was made but no response received
      throw {
        statusCode: 503,
        message: 'No response from server. Please try again later.'
      };
    } else {
      // Error setting up the request
      throw {
        statusCode: 500,
        message: 'Failed to send request: ' + error.message
      };
    }
  }
};

const EnhancedAnalysisForm = ({ onSubmit, isLoading }) => {
  const [url, setUrl] = useState('');
  const [error, setError] = useState('');

  const validateUrl = (value) => {
    const urlPattern = /^(https?:\/\/)?(www\.)?[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)$/;
    return urlPattern.test(value);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const trimmedUrl = url.trim();

    if (!trimmedUrl) {
      setError('URL is required');
      return;
    }

    if (!validateUrl(trimmedUrl)) {
      setError('Please enter a valid URL');
      return;
    }

    setError('');
    onSubmit(trimmedUrl);
  };

  return (
    <div className="mb-8 w-full max-w-2xl mx-auto bg-white rounded-lg shadow-md p-6 transform transition-all hover:shadow-lg">
      <h2 className="text-xl font-semibold mb-4 text-gray-800">Enter a URL to analyze</h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="url" className="block text-sm font-medium text-gray-700 mb-1">
            Web Page URL
          </label>
          <div className="relative mt-1 rounded-md shadow-sm">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M12.586 4.586a2 2 0 112.828 2.828l-3 3a2 2 0 01-2.828 0 1 1 0 00-1.414 1.414 4 4 0 005.656 0l3-3a4 4 0 00-5.656-5.656l-1.5 1.5a1 1 0 101.414 1.414l1.5-1.5zm-5 5a2 2 0 012.828 0 1 1 0 101.414-1.414 4 4 0 00-5.656 0l-3 3a4 4 0 105.656 5.656l1.5-1.5a1 1 0 10-1.414-1.414l-1.5 1.5a2 2 0 11-2.828-2.828l3-3z" clipRule="evenodd" />
              </svg>
            </div>
            <input
              type="text"
              id="url"
              name="url"
              placeholder="https://example.com"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              className={`block w-full pl-10 pr-12 py-3 border ${error ? 'border-red-500' : 'border-gray-300'
                } rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 text-gray-900`}
              disabled={isLoading}
            />
            {url && (
              <div className="absolute inset-y-0 right-0 pr-3 flex items-center">
                <button
                  type="button"
                  onClick={() => setUrl('')}
                  className="text-gray-400 hover:text-gray-500 focus:outline-none"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                  </svg>
                </button>
              </div>
            )}
          </div>
          {error && (
            <p className="mt-2 text-sm text-red-600">{error}</p>
          )}
        </div>

        <div className="flex justify-end">
          <button
            type="submit"
            disabled={isLoading}
            className="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
          >
            {isLoading ? (
              <>
                <svg className="animate-spin -ml-1 mr-2 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Analyzing...
              </>
            ) : (
              <span className="flex items-center">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clipRule="evenodd" />
                </svg>
                Analyze
              </span>
            )}
          </button>
        </div>
      </form>
    </div>
  );
};

const EnhancedAnalysisResult = ({ result }) => {
  if (!result) return null;

  const headings = {
    h1: 0, h2: 0, h3: 0, h4: 0, h5: 0, h6: 0,
    ...(result.headings || {})
  };

  const links = {
    internal: 0, external: 0, inaccessible: 0,
    ...(result.links || {})
  };

  const htmlVersion = result.htmlVersion || 'Unknown';
  const title = result.title || 'No title';
  const containsLoginForm = Boolean(result.containsLoginForm);

  const totalHeadings = Object.values(headings).reduce((a, b) => a + b, 0);
  const totalLinks = links.internal + links.external;

  return (
    <div className="mt-8 bg-white shadow overflow-hidden sm:rounded-lg">
      <div className="px-4 py-5 sm:px-6 bg-gradient-to-r from-blue-50 to-indigo-50">
        <h3 className="text-lg leading-6 font-medium text-gray-900 flex items-center">
          <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 mr-2 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Analysis Results
        </h3>
        <p className="mt-1 max-w-2xl text-sm text-gray-500">
          Detailed information about the analyzed web page.
        </p>
      </div>

      {/* Summary Stats */}
      <div className="border-t border-gray-200 px-4 py-5 sm:p-0">
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 p-4">
          <div className="bg-blue-50 p-4 rounded-lg text-center shadow-sm">
            <div className="text-3xl font-bold text-blue-600">{htmlVersion}</div>
            <div className="text-sm text-gray-500">HTML Version</div>
          </div>
          <div className="bg-green-50 p-4 rounded-lg text-center shadow-sm">
            <div className="text-3xl font-bold text-green-600">{totalHeadings}</div>
            <div className="text-sm text-gray-500">Total Headings</div>
          </div>
          <div className="bg-purple-50 p-4 rounded-lg text-center shadow-sm">
            <div className="text-3xl font-bold text-purple-600">{totalLinks}</div>
            <div className="text-sm text-gray-500">Total Links</div>
          </div>
        </div>
      </div>

      <div className="border-t border-gray-200">
        <dl>
          {/* HTML Version */}
          <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500 flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
              </svg>
              HTML Version
            </dt>
            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <span className="inline-flex items-center px-3 py-0.5 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
                {htmlVersion}
              </span>
            </dd>
          </div>

          {/* Page Title */}
          <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500 flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              Page Title
            </dt>
            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2 font-medium">
              {title}
            </dd>
          </div>

          {/* Headings */}
          <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500 flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h7" />
              </svg>
              Headings
            </dt>
            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <div className="grid grid-cols-6 gap-4">
                <div className="text-center p-2 rounded bg-green-50">
                  <span className="block text-lg font-semibold text-green-600">{headings.h1}</span>
                  <span className="text-xs text-gray-500">H1</span>
                </div>
                <div className="text-center p-2 rounded bg-green-50">
                  <span className="block text-lg font-semibold text-green-600">{headings.h2}</span>
                  <span className="text-xs text-gray-500">H2</span>
                </div>
                <div className="text-center p-2 rounded bg-green-50">
                  <span className="block text-lg font-semibold text-green-600">{headings.h3}</span>
                  <span className="text-xs text-gray-500">H3</span>
                </div>
                <div className="text-center p-2 rounded bg-green-50">
                  <span className="block text-lg font-semibold text-green-600">{headings.h4}</span>
                  <span className="text-xs text-gray-500">H4</span>
                </div>
                <div className="text-center p-2 rounded bg-green-50">
                  <span className="block text-lg font-semibold text-green-600">{headings.h5}</span>
                  <span className="text-xs text-gray-500">H5</span>
                </div>
                <div className="text-center p-2 rounded bg-green-50">
                  <span className="block text-lg font-semibold text-green-600">{headings.h6}</span>
                  <span className="text-xs text-gray-500">H6</span>
                </div>
              </div>
            </dd>
          </div>

          {/* Links */}
          <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500 flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
              </svg>
              Links
            </dt>
            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <div className="grid grid-cols-3 gap-4">
                <div className="text-center p-3 bg-blue-50 rounded-lg shadow-sm">
                  <span className="block text-lg font-semibold text-blue-600">{links.internal}</span>
                  <span className="text-xs text-gray-600">Internal</span>
                </div>
                <div className="text-center p-3 bg-green-50 rounded-lg shadow-sm">
                  <span className="block text-lg font-semibold text-green-600">{links.external}</span>
                  <span className="text-xs text-gray-600">External</span>
                </div>
                <div className="text-center p-3 bg-red-50 rounded-lg shadow-sm">
                  <span className="block text-lg font-semibold text-red-600">{links.inaccessible}</span>
                  <span className="text-xs text-gray-600">Inaccessible</span>
                </div>
              </div>
            </dd>
          </div>

          {/* Login Form */}
          <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500 flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
              Login Form
            </dt>
            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              {containsLoginForm ? (
                <span className="inline-flex items-center px-3 py-0.5 rounded-full text-sm font-medium bg-green-100 text-green-800">
                  <svg className="mr-1.5 h-2 w-2 text-green-600" fill="currentColor" viewBox="0 0 8 8">
                    <circle cx="4" cy="4" r="3" />
                  </svg>
                  Detected
                </span>
              ) : (
                <span className="inline-flex items-center px-3 py-0.5 rounded-full text-sm font-medium bg-gray-100 text-gray-800">
                  <svg className="mr-1.5 h-2 w-2 text-gray-400" fill="currentColor" viewBox="0 0 8 8">
                    <circle cx="4" cy="4" r="3" />
                  </svg>
                  Not detected
                </span>
              )}
            </dd>
          </div>
        </dl>
      </div>
    </div>
  );
};

const EnhancedErrorDisplay = ({ error }) => {
  if (!error) return null;

  return (
    <div className="bg-red-50 border-l-4 border-red-500 p-4 my-6 rounded-lg shadow-md">
      <div className="flex">
        <div className="flex-shrink-0">
          <svg className="h-6 w-6 text-red-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div className="ml-3">
          <h3 className="text-sm font-medium text-red-800">
            Error {error.statusCode && `(${error.statusCode})`}
          </h3>
          <div className="mt-2 text-sm text-red-700">
            <p>{error.message || 'An unknown error occurred'}</p>
          </div>
          <div className="mt-3">
            <div className="flex space-x-3">
              <button
                type="button"
                className="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded text-red-700 bg-red-100 hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
              >
                Try Again
              </button>
              <a
                href="https://developer.mozilla.org/en-US/docs/Web/HTTP/Status"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center px-3 py-1.5 border border-transparent text-xs font-medium rounded text-gray-700 bg-gray-100 hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
              >
                Learn About HTTP Status Codes
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

const EnhancedLoadingIndicator = () => {
  return (
    <div className="flex flex-col justify-center items-center my-12 p-8 bg-white rounded-lg shadow-md">
      <div className="animate-spin rounded-full h-16 w-16 border-t-4 border-b-4 border-blue-500"></div>
      <span className="mt-4 text-lg text-gray-700">Analyzing web page...</span>
      <p className="mt-2 text-sm text-gray-500">This may take a moment depending on the site size</p>
    </div>
  );
};

const EnhancedRecentAnalysis = ({ url, onClick }) => {
  if (!url) return null;
  return (
    <div className="flex items-center bg-blue-50 p-2 rounded-md mb-4 text-sm">
      <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-blue-500 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span className="text-gray-600 truncate flex-grow">Recently analyzed:</span>
      <button
        onClick={() => onClick(url)}
        className="ml-2 text-blue-600 hover:text-blue-800 truncate max-w-xs"
      >
        {url}
      </button>
    </div>
  );
};

function App() {
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);
  const [analyzedUrl, setAnalyzedUrl] = useState('');
  const [recentUrls, setRecentUrls] = useState([]);

  const handleAnalyzeUrl = async (url) => {
    setIsLoading(true);
    setError(null);
    setResult(null);

    try {
      const data = await analyzeWebPage(url);
      const safeData = {
        htmlVersion: data?.htmlVersion || 'Unknown',
        title: data?.title || 'No title',
        headings: {
          h1: 0, h2: 0, h3: 0, h4: 0, h5: 0, h6: 0,
          ...(data?.headings || {})
        },
        links: {
          internal: 0, external: 0, inaccessible: 0,
          ...(data?.links || {})
        },
        containsLoginForm: Boolean(data?.containsLoginForm)
      };

      setResult(safeData);
      setAnalyzedUrl(url);

      // Add to recent URLs if not already present
      if (!recentUrls.includes(url)) {
        setRecentUrls(prev => [url, ...prev].slice(0, 5));
      }
    } catch (err) {
      setError(err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 text-blue-600 mr-3" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z" clipRule="evenodd" />
              </svg>
              <h1 className="text-3xl font-bold text-gray-900">Web Page Analyzer</h1>
            </div>
            <div className="text-sm text-gray-500">
              <span>Analyze webpages instantly</span>
            </div>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div className="bg-white shadow overflow-hidden sm:rounded-lg mb-6 p-6">
          <div className="flex items-start">
            <div className="flex-shrink-0 pt-1">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div className="ml-3">
              <p className="text-gray-700">
                Enter a URL below to analyze the HTML structure, links, and more. The analyzer will check the page's HTML version, count headings, categorize links, and detect login forms.
              </p>
            </div>
          </div>
        </div>

        {recentUrls.length > 0 && (
          <div className="mb-6">
            <h3 className="text-sm font-medium text-gray-500 mb-2">Recent analyses</h3>
            <div className="space-y-2">
              {recentUrls.map((url, index) => (
                <EnhancedRecentAnalysis
                  key={index}
                  url={url}
                  onClick={handleAnalyzeUrl}
                />
              ))}
            </div>
          </div>
        )}

        <EnhancedAnalysisForm onSubmit={handleAnalyzeUrl} isLoading={isLoading} />

        <EnhancedErrorDisplay error={error} />

        {isLoading ? (
          <EnhancedLoadingIndicator />
        ) : (
          result && (
            <>
              {analyzedUrl && (
                <div className="mb-4 text-sm bg-gray-100 p-3 rounded-md">
                  <span className="text-gray-500">Analysis for: </span>
                  <a href={analyzedUrl} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline font-medium">
                    {analyzedUrl}
                  </a>
                </div>
              )}
              <EnhancedAnalysisResult result={result} />
            </>
          )
        )}
      </main>
    </div>
  );
}
export default App;
