"use client";

import { useState } from 'react';
import { motion } from 'framer-motion';
import { Project } from '@/entities/Project';

type Category = 'all' | Project['category'];

interface ProjectFilterProps {
  onFilterChange: (category: Category) => void;
  activeCategory: Category;
}

export const ProjectFilter = ({
  onFilterChange,
  activeCategory,
}: ProjectFilterProps) => {
  const categories: { value: Category; label: string }[] = [
    { value: 'all', label: 'Все проекты' },
    { value: 'frontend', label: 'Frontend' },
    { value: 'backend', label: 'Backend' },
    { value: 'fullstack', label: 'Fullstack' },
    { value: 'devops', label: 'DevOps' },
    { value: 'other', label: 'Другое' },
  ];

  return (
    <div className="flex flex-wrap gap-2 justify-center mb-10">
      {categories.map((category) => (
        <motion.button
          key={category.value}
          onClick={() => onFilterChange(category.value)}
          className={`px-4 py-2 rounded-full transition-colors ${
            activeCategory === category.value
              ? 'bg-blue-600 text-white'
              : 'bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700'
          }`}
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
        >
          {category.label}
        </motion.button>
      ))}
    </div>
  );
};
