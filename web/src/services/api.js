/* eslint-disable no-throw-literal */
import axios from 'axios';


const API_URL = process.env.REACT_APP_API_URL || '/api';

/**
 * Analyze web page URL
 * @param {string} url - url need send/analys
 * @returns
 */
export const analyzeWebPage = async (url) => {
  try {
    const response = await axios.post(`${API_URL}/analyze`, { url });
    return response.data;
  } catch (error) {
    if (error.response) {
      throw error.response.data;
    } else if (error.request) {
      throw {
        statusCode: 503,
        message: 'No response got from server. Please try again later.'
      };
    } else {
      throw {
        statusCode: 500,
        message: 'Failed to send request: ' + error.message
      };
    }
  }
};