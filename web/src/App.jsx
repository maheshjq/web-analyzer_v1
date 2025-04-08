import React, { useState } from 'react';
import AnalysisForm from './components/AnalysisForm';
import AnalysisResult from './components/AnalysisResult';
import LoadingIndicator from './components/LoadingIndicator';
import ErrorDisplay from './components/ErrorDisplay';
import { analyzeWebPage } from './services/api';

function App() {
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);
  const [analyzedUrl, setAnalyzedUrl] = useState('');

  const handleAnalyzeUrl = async (url) => {
    setIsLoading(true);
    setError(null);
    setResult(null); // Clear previous results
    
    try {
      // Analyze the URL
      const data = await analyzeWebPage(url);
      // Ensure data has all required structures before setting state
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
    } catch (err) {
      setError(err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
          <h1 className="text-3xl font-bold text-gray-900">Web Page Analyzer</h1>
        </div>
      </header>
      
      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div className="bg-white shadow overflow-hidden sm:rounded-lg mb-6 p-6">
          <p className="text-gray-700">
            Enter a URL below to analyze the HTML structure, links, and more.
          </p>
        </div>
        
        <AnalysisForm onSubmit={handleAnalyzeUrl} isLoading={isLoading} />
        
        <ErrorDisplay error={error} />
        
        {isLoading ? (
          <LoadingIndicator />
        ) : (
          result && (
            <>
              {analyzedUrl && (
                <div className="mb-4 text-sm text-gray-600">
                  Analysis for: <a href={analyzedUrl} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">{analyzedUrl}</a>
                </div>
              )}
              <AnalysisResult result={result} />
            </>
          )
        )}
      </main>
      
      <footer className="bg-white border-t border-gray-200 py-6">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <p className="text-center text-gray-500 text-sm">
            &copy; {new Date().getFullYear()} Web Page Analyzer By Mahesh
          </p>
        </div>
      </footer>
    </div>
  );
}

export default App;