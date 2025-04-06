import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || '/api';

/**
 * Analyze a web page URL
 * @param {string} url - The URL to analyze
 * @returns {Promise} - Promise with analysis results
 */
export const analyzeWebPage = async (url) => {
  try {
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