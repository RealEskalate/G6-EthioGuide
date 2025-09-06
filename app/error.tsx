"use client";

import Link from 'next/link';
import React from 'react';

const ErrorPage = () => {
  return (
    <div className="bg-gray-50 flex flex-col items-center justify-center min-h-screen p-4 font-sans text-gray-800">
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');
        .font-sans {
          font-family: 'Inter', sans-serif;
        }
      `}</style>
      <div className="text-center p-8 bg-white rounded-xl shadow-lg max-w-lg w-full">
        <div className="mb-4">
          <svg className="w-24 h-24 mx-auto text-rose-600 mb-2" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 15h2v2h-2v-2zm0-10h2v8h-2V7z"/>
          </svg>
          <div className="text-8xl font-extrabold text-rose-600">500</div>
        </div>
        <h1 className="text-3xl sm:text-4xl font-bold mb-2">
          Something went wrong
        </h1>
        <p className="text-base sm:text-lg text-gray-600 mb-6">
          We&apos;re sorry, but an unexpected error occurred. Please try again later.
        </p>
        <div className="flex flex-col sm:flex-row justify-center items-center gap-4">
          <Link
            href="/"
            className="w-full sm:w-auto px-6 py-3 text-sm sm:text-base font-semibold text-white bg-rose-600 rounded-lg hover:bg-rose-700 transition-colors shadow-md"
          >
            Go Home
          </Link>
          <a
            href="#"
            onClick={() => window.history.back()}
            className="w-full sm:w-auto px-6 py-3 text-sm sm:text-base font-semibold text-rose-600 bg-transparent border-2 border-rose-600 rounded-lg hover:bg-rose-50 transition-colors"
          >
            Go Back
          </a>
        </div>
      </div>
    </div>
  );
};

export default ErrorPage;
