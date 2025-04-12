import React from 'react';

const AnalysisVisualization = ({ result }) => {
  if (!result) return null;
  
  const headings = result.headings || { h1: 0, h2: 0, h3: 0, h4: 0, h5: 0, h6: 0 };
  const links = result.links || { internal: 0, external: 0, inaccessible: 0 };
  
  // Calculate maximum values for scaling
  const maxHeading = Math.max(...Object.values(headings), 1);
  const maxLink = Math.max(links.internal, links.external, links.inaccessible, 1);
  
  // Generate heading bars
  const headingBars = [
    { type: 'H1', count: headings.h1, color: 'bg-green-500' },
    { type: 'H2', count: headings.h2, color: 'bg-green-400' },
    { type: 'H3', count: headings.h3, color: 'bg-green-300' },
    { type: 'H4', count: headings.h4, color: 'bg-blue-300' },
    { type: 'H5', count: headings.h5, color: 'bg-blue-200' },
    { type: 'H6', count: headings.h6, color: 'bg-blue-100' },
  ];
  
  // Generate link visualization
  const linkData = [
    { type: 'Internal', count: links.internal, color: 'bg-blue-500' },
    { type: 'External', count: links.external, color: 'bg-green-500' },
    { type: 'Inaccessible', count: links.inaccessible, color: 'bg-red-500' },
  ];

  return (
    <div className="mt-8 bg-white shadow overflow-hidden sm:rounded-lg">
      <div className="px-4 py-5 sm:px-6 bg-gradient-to-r from-blue-50 to-purple-50">
        <h3 className="text-lg leading-6 font-medium text-gray-900 flex items-center">
          <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 mr-2 text-purple-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          Page Structure Visualization
        </h3>
        <p className="mt-1 max-w-2xl text-sm text-gray-500">
          Visual representation of headings and links
        </p>
      </div>

      <div className="border-t border-gray-200 px-4 py-5">
        <h4 className="text-md font-medium text-gray-700 mb-4">Heading Distribution</h4>
        <div className="space-y-3">
          {headingBars.map((heading) => (
            <div key={heading.type} className="flex items-center">
              <div className="w-10 text-right mr-3 text-sm font-medium text-gray-600">
                {heading.type}
              </div>
              <div className="relative flex-grow h-8 bg-gray-100 rounded-r">
                <div 
                  className={`absolute top-0 left-0 h-full ${heading.color} rounded-r transition-all duration-500 ease-in-out`} 
                  style={{ width: `${(heading.count / maxHeading) * 100}%` }}
                ></div>
                <div className="absolute inset-0 flex items-center justify-end pr-3">
                  <span className="text-sm font-medium text-gray-700">{heading.count}</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="border-t border-gray-200 px-4 py-5">
        <h4 className="text-md font-medium text-gray-700 mb-4">Link Analysis</h4>
        
        {/* Pie chart representation using CSS */}
        <div className="flex items-center justify-around">
          <div className="relative w-40 h-40">
            {links.internal + links.external + links.inaccessible > 0 ? (
              <svg viewBox="0 0 100 100" className="w-full h-full">
                {/* Calculate and draw sectors */}
                {(() => {
                  const total = links.internal + links.external + links.inaccessible;
                  let startAngle = 0;
                  
                  return linkData
                    .filter(item => item.count > 0)
                    .map((item, index) => {
                      const percentage = item.count / total;
                      const endAngle = startAngle + percentage * 360;
                      
                      // Convert to radians and calculate x,y coordinates
                      const startRad = (startAngle - 90) * Math.PI / 180;
                      const endRad = (endAngle - 90) * Math.PI / 180;
                      
                      const x1 = 50 + 50 * Math.cos(startRad);
                      const y1 = 50 + 50 * Math.sin(startRad);
                      const x2 = 50 + 50 * Math.cos(endRad);
                      const y2 = 50 + 50 * Math.sin(endRad);
                      
                      // Determine which arc to draw (large or small)
                      const largeArcFlag = percentage > 0.5 ? 1 : 0;
                      
                      // Create path
                      const pathData = [
                        `M 50 50`,
                        `L ${x1} ${y1}`,
                        `A 50 50 0 ${largeArcFlag} 1 ${x2} ${y2}`,
                        'Z'
                      ].join(' ');
                      
                      // Update start angle for next segment
                      const result = (
                        <path 
                          key={index}
                          d={pathData}
                          className={item.color}
                        />
                      );
                      
                      startAngle = endAngle;
                      return result;
                    });
                })()}
              </svg>
            ) : (
              <div className="w-full h-full flex items-center justify-center bg-gray-100 rounded-full">
                <span className="text-gray-500 text-sm">No links</span>
              </div>
            )}
          </div>
          
          <div className="space-y-2">
            {linkData.map((item) => (
              <div key={item.type} className="flex items-center">
                <div className={`w-4 h-4 rounded ${item.color} mr-2`}></div>
                <span className="text-sm text-gray-700">{item.type}: {item.count}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default AnalysisVisualization;