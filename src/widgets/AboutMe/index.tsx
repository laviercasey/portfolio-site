"use client";

import { motion, useInView, AnimatePresence } from 'framer-motion';
import Image from 'next/image';
import { useRef, useState } from 'react';
import { profile } from '@/entities/Profile';
import { getAssetPath } from '@/shared/lib/utils';

export const AboutMe = () => {
  const [showAllSkills, setShowAllSkills] = useState(false);
  
  const sectionRef = useRef(null);
  const skillsRef = useRef(null);
  const experienceRef = useRef(null);
  
  const isSectionInView = useInView(sectionRef, { once: true, amount: 0.2 });
  const isSkillsInView = useInView(skillsRef, { once: true, amount: 0.5 });
  const isExperienceInView = useInView(experienceRef, { once: true, amount: 0.3 });

  const visibleSkills = showAllSkills ? profile.skills : profile.skills.slice(0, 6);

  const titleVariants = {
    hidden: { opacity: 0, y: -20 },
    visible: { 
      opacity: 1, 
      y: 0,
      transition: { duration: 0.6 }
    }
  };

  const leftBlockVariants = {
    hidden: { opacity: 0, x: -100 },
    visible: { 
      opacity: 1, 
      x: 0,
      transition: { 
        duration: 0.8,
        delay: 0.3
      }
    }
  };

  const rightBlockVariants = {
    hidden: { opacity: 0, x: 100 },
    visible: { 
      opacity: 1, 
      x: 0,
      transition: { 
        duration: 0.8,
        delay: 0.3
      }
    }
  };

  const skillsContainerVariants = {
    hidden: { opacity: 0 },
    visible: { 
      opacity: 1,
      transition: { 
        staggerChildren: 0.1,
        delayChildren: 0.2,
        when: "beforeChildren" 
      }
    }
  };

  const skillItemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { 
      opacity: 1, 
      y: 0,
      transition: { duration: 0.5 }
    }
  };

  const buttonVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { 
      opacity: 1, 
      y: 0,
      transition: { duration: 0.5, delay: 0.7 }
    }
  };

  const experienceContainerVariants = {
    hidden: { opacity: 0 },
    visible: { 
      opacity: 1,
      transition: { 
        staggerChildren: 0.2,
        delayChildren: 0.3,
        when: "beforeChildren" 
      }
    }
  };

  const experienceItemVariants = {
    hidden: { opacity: 0, y: -30 },
    visible: { 
      opacity: 1, 
      y: 0,
      transition: { duration: 0.6 }
    }
  };

  return (
    <section ref={sectionRef} className="py-20 bg-gray-50 dark:bg-gray-900">
      <div className="container mx-auto px-4 md:px-6">
        <motion.div
          className="max-w-4xl mx-auto text-center mb-16"
          variants={titleVariants}
          initial="hidden"
          animate={isSectionInView ? "visible" : "hidden"}
        >
          <h2 className="text-3xl md:text-4xl font-bold mb-4">Обо мне</h2>
          <p className="text-xl text-gray-600 dark:text-gray-300">
            {profile.bio.short}
          </p>
        </motion.div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-start">
          <motion.div
            variants={leftBlockVariants}
            initial="hidden"
            animate={isSectionInView ? "visible" : "hidden"}
          >
            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl overflow-hidden">
              <div className="relative h-64 w-full">
                <Image
                  src={getAssetPath("/images/about-me.jpg")}
                  alt={`${profile.name} за работой`}
                  fill
                  className="object-cover"
                  sizes="(max-width: 768px) 100vw, 600px"
                />
              </div>
              <div className="p-6">
                <h3 className="text-2xl font-bold mb-4">Мой путь в разработке</h3>
                <p className="text-gray-600 dark:text-gray-300 mb-6">
                  {profile.bio.full}
                </p>
              </div>
            </div>
          </motion.div>

          <div>
            <motion.div
              ref={skillsRef}
              className="bg-white dark:bg-gray-800 rounded-xl shadow-xl p-6 mb-8"
              variants={rightBlockVariants}
              initial="hidden"
              animate={isSectionInView ? "visible" : "hidden"}
            >
              <h3 className="text-2xl font-bold mb-6">Навыки</h3>
              
              <motion.div 
                className={`space-y-4 ${showAllSkills ? 'max-h-96 overflow-y-auto pr-2' : ''}`}
                variants={skillsContainerVariants}
                initial="hidden"
                animate={isSkillsInView ? "visible" : "hidden"}
                style={{
                  scrollbarWidth: 'thin',
                  scrollbarColor: '#3b82f6 #e5e7eb'
                }}
              >
                <AnimatePresence>
                  {visibleSkills.map((skill, index) => (
                    <motion.div 
                      key={skill.name} 
                      variants={skillItemVariants}
                      initial={index >= 6 ? { opacity: 0, y: 20 } : undefined}
                      animate={{ opacity: 1, y: 0 }}
                      exit={{ opacity: 0, y: 20 }}
                      transition={{ duration: 0.5, delay: index >= 6 ? index * 0.05 : 0 }}
                    >
                      <div className="flex justify-between mb-1">
                        <span className="font-medium">{skill.name}</span>
                        <span className="text-gray-500 dark:text-gray-400">{skill.level}%</span>
                      </div>
                      <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
                        <motion.div
                          className="bg-blue-600 dark:bg-blue-500 h-2.5 rounded-full"
                          initial={{ width: 0 }}
                          animate={isSkillsInView ? { width: `${skill.level}%` } : { width: 0 }}
                          transition={{ duration: 1, delay: index >= 6 ? 0.2 : 0.5 }}
                        ></motion.div>
                      </div>
                    </motion.div>
                  ))}
                </AnimatePresence>
              </motion.div>
              
              {profile.skills.length > 6 && (
                <motion.div
                  className="mt-6 text-center"
                  variants={buttonVariants}
                  initial="hidden"
                  animate={isSkillsInView ? "visible" : "hidden"}
                >
                  <motion.button
                    className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition-colors duration-300"
                    onClick={() => setShowAllSkills(!showAllSkills)}
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                  >
                    {showAllSkills ? 'Показать меньше' : 'Показать больше'}
                  </motion.button>
                </motion.div>
              )}
            </motion.div>

            <motion.div
              ref={experienceRef}
              className="bg-white dark:bg-gray-800 rounded-xl shadow-xl p-6"
              variants={rightBlockVariants}
              initial="hidden"
              animate={isSectionInView ? "visible" : "hidden"}
            >
              <h3 className="text-2xl font-bold mb-6">Опыт работы</h3>
              <motion.div 
                className="space-y-6"
                variants={experienceContainerVariants}
                initial="hidden"
                animate={isExperienceInView ? "visible" : "hidden"}
              >
                {profile.experience.map((exp, index) => (
                  <motion.div
                    key={index}
                    variants={experienceItemVariants}
                    className="relative pl-8 pb-6 border-l-2 border-gray-200 dark:border-gray-700 last:border-0 last:pb-0"
                  >
                    <div className="absolute left-[-9px] top-0 w-4 h-4 bg-blue-600 dark:bg-blue-500 rounded-full"></div>
                    <h4 className="text-xl font-bold">{exp.position}</h4>
                    <p className="text-blue-600 dark:text-blue-400 mb-1">{exp.company}</p>
                    <p className="text-gray-500 dark:text-gray-400 text-sm mb-2">{exp.period}</p>
                    <p className="text-gray-600 dark:text-gray-300">{exp.description}</p>
                  </motion.div>
                ))}
              </motion.div>
            </motion.div>
          </div>
        </div>
      </div>
    </section>
  );
};
