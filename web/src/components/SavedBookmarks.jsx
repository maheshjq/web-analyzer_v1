import React, { useState, useEffect } from 'react';

const SavedBookmarks = ({ onSelectUrl }) => {
  const [bookmarks, setBookmarks] = useState([]);
  const [newBookmark, setNewBookmark] = useState('');
  const [isAdding, setIsAdding] = useState(false);

  // Load bookmarks from localStorage on component mount
  useEffect(() => {
    const savedBookmarks = localStorage.getItem('webAnalyzerBookmarks');
    if (savedBookmarks) {
      try {
        setBookmarks(JSON.parse(savedBookmarks));
      } catch (e) {
        console.error('Error parsing bookmarks', e);
      }
    }
  }, []);

  // Save bookmarks to localStorage whenever they change
  useEffect(() => {
    localStorage.setItem('webAnalyzerBookmarks', JSON.stringify(bookmarks));
  }, [bookmarks]);

  const addBookmark = () => {
    if (newBookmark && !bookmarks.includes(newBookmark)) {
      setBookmarks([...bookmarks, newBookmark]);
      setNewBookmark('');
      setIsAdding(false);
    }
  };

  const removeBookmark = (url) => {
    setBookmarks(bookmarks.filter(bookmark => bookmark !== url));
  };

  return (
    <div className="mt-6 mb-8">
      <div className="flex justify-between items-center mb-3">
        <h3 className="text-sm font-medium text-gray-500">Saved Sites</h3>
        {!isAdding && (
          <button
            onClick={() => setIsAdding(true)}
            className="text-sm text-blue-600 hover:text-blue-800"
          >
            + Add New
          </button>
        )}
      </div>

      {isAdding && (
        <div className="mb-4 flex space-x-2">
          <input
            type="text"
            value={newBookmark}
            onChange={(e) => setNewBookmark(e.target.value)}
            placeholder="https://example.com"
            className="flex-grow border border-gray-300 rounded-md shadow-sm px-3 py-2 text-sm focus:ring-blue-500 focus:border-blue-500"
          />
          <button
            onClick={addBookmark}
            className="inline-flex items-center px-3 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            Save
          </button>
          <button
            onClick={() => {
              setIsAdding(false);
              setNewBookmark('');
            }}
            className="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            Cancel
          </button>
        </div>
      )}

      {bookmarks.length > 0 ? (
        <div className="bg-white shadow overflow-hidden sm:rounded-lg divide-y divide-gray-200">
          {bookmarks.map((bookmark, index) => (
            <div key={index} className="px-4 py-3 flex justify-between items-center hover:bg-gray-50">
              <button
                onClick={() => onSelectUrl(bookmark)}
                className="text-sm text-blue-600 hover:text-blue-800 truncate max-w-xs text-left"
              >
                {bookmark}
              </button>
              <button
                onClick={() => removeBookmark(bookmark)}
                className="text-gray-400 hover:text-red-500"
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          ))}
        </div>
      ) : (
        <div className="text-center py-6 bg-white rounded-lg shadow-sm">
          <p className="text-gray-500 text-sm">No saved sites yet</p>
        </div>
      )}
    </div>
  );
};

export default SavedBookmarks;