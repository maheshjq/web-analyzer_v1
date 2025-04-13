import React from 'react';

const AnalysisResult = (props) => {
  if (!props || !props.result) {
    return null;
  }

  const result = props.result;
  
  const headings = {
    h1: 0,
    h2: 0,
    h3: 0,
    h4: 0,
    h5: 0,
    h6: 0,
    ...(result.headings || {})
  };

  const links = {
    internal: 0,
    external: 0,
    inaccessible: 0,
    ...(result.links || {})
  };

  const htmlVersion = result.htmlVersion || 'Unknown';
  const title = result.title || 'No title';
  const containsLoginForm = Boolean(result.containsLoginForm);

  return (
    <div className="mt-8 bg-white shadow overflow-hidden sm:rounded-lg">
      <div className="px-4 py-5 sm:px-6">
        <h3 className="text-lg leading-6 font-medium text-gray-900">
          Analysis Results
        </h3>
        <p className="mt-1 max-w-2xl text-sm text-gray-500">
          Detailed information about the analyzed web page.
        </p>
      </div>

      <div className="border-t border-gray-200">
        <dl>
          {/* HTML Version */}
          <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500">HTML Version</dt>
            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              {htmlVersion}
            </dd>
          </div>

          {/* Page Title */}
          <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500">Page Title</dt>
            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              {title}
            </dd>
          </div>

          {/* Headings */}
          <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500">Headings</dt>
            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <div className="grid grid-cols-6 gap-4">
                <div className="text-center">
                  <span className="block text-lg font-semibold">{headings.h1}</span>
                  <span className="text-xs text-gray-500">H1</span>
                </div>
                <div className="text-center">
                  <span className="block text-lg font-semibold">{headings.h2}</span>
                  <span className="text-xs text-gray-500">H2</span>
                </div>
                <div className="text-center">
                  <span className="block text-lg font-semibold">{headings.h3}</span>
                  <span className="text-xs text-gray-500">H3</span>
                </div>
                <div className="text-center">
                  <span className="block text-lg font-semibold">{headings.h4}</span>
                  <span className="text-xs text-gray-500">H4</span>
                </div>
                <div className="text-center">
                  <span className="block text-lg font-semibold">{headings.h5}</span>
                  <span className="text-xs text-gray-500">H5</span>
                </div>
                <div className="text-center">
                  <span className="block text-lg font-semibold">{headings.h6}</span>
                  <span className="text-xs text-gray-500">H6</span>
                </div>
              </div>
            </dd>
          </div>

          {/* Links */}
          <div className="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500">Links</dt>
            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              <div className="grid grid-cols-3 gap-4">
                <div className="text-center p-3 bg-blue-50 rounded">
                  <span className="block text-lg font-semibold">{links.internal}</span>
                  <span className="text-xs text-gray-600">Internal</span>
                </div>
                <div className="text-center p-3 bg-green-50 rounded">
                  <span className="block text-lg font-semibold">{links.external}</span>
                  <span className="text-xs text-gray-600">External</span>
                </div>
                <div className="text-center p-3 bg-red-50 rounded">
                  <span className="block text-lg font-semibold">{links.inaccessible}</span>
                  <span className="text-xs text-gray-600">Inaccessible</span>
                </div>
              </div>
            </dd>
          </div>

          {/* Login Form */}
          <div className="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-500">Login Form</dt>
            <dd className="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              {containsLoginForm ? (
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                  Detected
                </span>
              ) : (
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
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

export default AnalysisResult;