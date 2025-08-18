"use client";

import { useState, useEffect, useRef } from 'react';
import { motion, useInView, AnimatePresence } from 'framer-motion';
import { ProjectCard, Project, projects } from '@/entities/Project';
import { ProjectFilter } from '@/features/ProjectFilter';

type Category = 'all' | Project['category'];

export const ProjectList = () => {
  const [activeCategory, setActiveCategory] = useState<Category>('all');
  const [filteredProjects, setFilteredProjects] = useState<Project[]>([]);
  const sectionRef = useRef(null);
  const isSectionInView = useInView(sectionRef, { once: true, amount: 0.2 });

  useEffect(() => {
    if (activeCategory === 'all') {
      setFilteredProjects(projects);
    } else {
      setFilteredProjects(projects.filter(project => project.category === activeCategory));
    }
  }, [activeCategory]);

  const getAnimationDirection = (index: number) => {
    const col = index % 3;
    const row = Math.floor(index / 3);
    
    if (row === 0) {
      if (col === 0) return { x: -50, y: 0 };
      if (col === 1) return { x: 0, y: -50 };
      if (col === 2) return { x: 50, y: 0 };
    }
    
    if (row === 1) {
      if (col === 0) return { x: -50, y: 0 };
      if (col === 1) return { x: 0, y: 50 };
      if (col === 2) return { x: 50, y: 0 };
    }
    
    return row % 2 === 0 
      ? { x: col % 2 === 0 ? -50 : 50, y: 0 } 
      : { x: 0, y: col % 2 === 0 ? -50 : 50 }; 
  };

  const titleVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { 
      opacity: 1, 
      y: 0,
      transition: { duration: 0.6 }
    }
  };

  const filterVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { 
      opacity: 1, 
      y: 0,
      transition: { 
        duration: 0.6,
        delay: 0.2
      }
    }
  };

  return (
    <section ref={sectionRef} className="py-20">
      <div className="container mx-auto px-4 md:px-6">
        <motion.div
          className="text-center mb-12"
          variants={titleVariants}
          initial="hidden"
          animate={isSectionInView ? "visible" : "hidden"}
        >
          <h2 className="text-3xl md:text-4xl font-bold mb-4">Мои проекты</h2>
          <p className="text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
            Исследуйте мои последние работы, демонстрирующие мой опыт в разработке фронтенда, бэкенда и полнофункциональных приложений.
          </p>
        </motion.div>

        <motion.div
          variants={filterVariants}
          initial="hidden"
          animate={isSectionInView ? "visible" : "hidden"}
        >
          <ProjectFilter
            onFilterChange={setActiveCategory}
            activeCategory={activeCategory}
          />
        </motion.div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <AnimatePresence mode="wait">
            {filteredProjects.map((project, index) => {
              const direction = getAnimationDirection(index);
              
              return (
                <motion.div
                  key={project.id}
                  layout
                  initial={{ 
                    opacity: 0, 
                    x: direction.x, 
                    y: direction.y 
                  }}
                  animate={{ 
                    opacity: 1, 
                    x: 0, 
                    y: 0,
                    transition: {
                      type: "spring",
                      damping: 15,
                      stiffness: 100,
                      delay: 0.3 + (index * 0.1)
                    }
                  }}
                  exit={{ 
                    opacity: 0, 
                    scale: 0.8,
                    transition: { duration: 0.3 }
                  }}
                  whileHover={{ 
                    y: -10,
                    transition: { duration: 0.2 }
                  }}
                >
                  <ProjectCard project={project} />
                </motion.div>
              );
            })}
          </AnimatePresence>
        </div>

        {filteredProjects.length === 0 && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="text-center py-12"
          >
            <p className="text-gray-500 dark:text-gray-400 text-lg">
              Проекты в этой категории отсутствуют
            </p>
          </motion.div>
        )}
      </div>
    </section>
  );
};
