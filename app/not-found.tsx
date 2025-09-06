"use client";

import React from 'react';

const NotFoundPage = () => {
  return (
    <div className="bg-background-light flex flex-col items-center justify-center min-h-screen p-4 font-sans text-neutral-dark">
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');
        .font-sans {
          font-family: 'Inter', sans-serif;
        }
      `}</style>
      <div className="text-center p-8 bg-white rounded-xl shadow-lg max-w-lg w-full border-neutral-light">
        <h1 className="text-6xl sm:text-8xl font-extrabold text-primary mb-4">
          404
        </h1>
        <h2 className="text-2xl sm:text-3xl font-bold mb-2">
          Page Not Found
        </h2>
        <p className="text-base sm:text-lg text-neutral-dark mb-6">
          The page you are looking for does not exist or has been moved.
        </p>
        <div className="flex flex-col sm:flex-row justify-center items-center gap-4">
          <button
            onClick={() => window.history.back()}
            className="w-full sm:w-auto px-6 py-3 text-sm sm:text-base font-semibold text-primary bg-transparent border-2 border-primary rounded-lg hover:bg-primary hover:text-white transition-colors"
          >
            Go Back
          </button>
         
        </div>
      </div>
    </div>
  );
};

export default NotFoundPage;
