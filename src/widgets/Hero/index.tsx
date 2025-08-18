"use client";

import { useState, useEffect } from 'react';
import { motion, useAnimation, useReducedMotion } from 'framer-motion';
import Image from 'next/image';
import { Button } from '@/shared/ui/Button';
import { CollaborationModal } from '@/features/CollaborationModal';
import { profile } from '@/entities/Profile';
import { getAssetPath } from '@/shared/lib/utils';

export const Hero = () => {
const [isModalOpen, setIsModalOpen] = useState(false);
const controls = useAnimation();
const shouldReduceMotion = useReducedMotion();

const techIcons = [
{ name: 'React', src: getAssetPath('/images/icons/react.svg'), size: 80 },
{ name: 'TypeScript', src: getAssetPath('/images/icons/typescript.svg'), size: 80 },
{ name: 'Next.js', src: getAssetPath('/images/icons/nextjs.svg'), size: 80 },
{ name: 'Vue.js', src: getAssetPath('/images/icons/vue.svg'), size: 80 },
{ name: 'Tailwind', src: getAssetPath('/images/icons/tailwind.svg'), size: 80 },
{ name: 'Postgresql', src: getAssetPath('/images/icons/postgresql.svg'), size: 80 },
];

useEffect(() => {
controls.start('visible');
}, [controls]);

const backgroundVariants = {
hidden: { opacity: 0 },
visible: { opacity: 0.8, transition: { duration: shouldReduceMotion ? 0.1 : 1.5 } }
};

const contentVariants = {
hidden: { opacity: 0, y: shouldReduceMotion ? 0 : 50 },
visible: { 
    opacity: 1, 
    y: 0,
    transition: { 
    duration: shouldReduceMotion ? 0.1 : 0.8,
    delay: shouldReduceMotion ? 0 : 0.3
    }
}
};

const avatarVariants = {
hidden: { opacity: 0, scale: shouldReduceMotion ? 1 : 0.8 },
visible: { 
    opacity: 1, 
    scale: 1,
    transition: { 
    duration: shouldReduceMotion ? 0.1 : 0.8,
    delay: shouldReduceMotion ? 0 : 0.5
    }
}
};

return (
<section className="relative min-h-screen flex items-center overflow-hidden py-5 md:py-0">
    <motion.div
    className="absolute inset-0 z-0"
    initial="hidden"
    animate="visible"
    variants={backgroundVariants}
    >
    <div className="absolute inset-0 bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-blue-900"></div>
    
    <motion.div 
        className="absolute top-1/4 left-1/4 w-96 h-96 rounded-full bg-blue-200 dark:bg-blue-800 opacity-30 blur-3xl"
        animate={shouldReduceMotion ? {} : { 
        scale: [1, 1.2, 1],
        x: [0, -30, 0],
        y: [0, 20, 0],
        }}
        transition={shouldReduceMotion ? {} : { 
        duration: 15,
        repeat: Infinity,
        repeatType: "reverse"
        }}
        style={{ willChange: shouldReduceMotion ? 'auto' : 'transform' }}
    />
    
    <motion.div 
        className="absolute bottom-1/3 right-1/3 w-64 h-64 rounded-full bg-purple-200 dark:bg-purple-900 opacity-20 blur-3xl"
        animate={shouldReduceMotion ? {} : { 
        scale: [1, 1.3, 1],
        x: [0, 40, 0],
        y: [0, -30, 0],
        }}
        transition={shouldReduceMotion ? {} : { 
        duration: 18,
        repeat: Infinity,
        repeatType: "reverse",
        delay: 2
        }}
        style={{ willChange: shouldReduceMotion ? 'auto' : 'transform' }}
    />
    
    <motion.div 
        className="absolute top-2/3 right-1/4 w-80 h-80 rounded-full bg-indigo-200 dark:bg-indigo-800 opacity-25 blur-3xl"
        animate={shouldReduceMotion ? {} : { 
        scale: [1, 1.1, 1],
        x: [0, 20, 0],
        y: [0, 40, 0],
        }}
        transition={shouldReduceMotion ? {} : { 
        duration: 20,
        repeat: Infinity,
        repeatType: "reverse",
        delay: 1
        }}
        style={{ willChange: shouldReduceMotion ? 'auto' : 'transform' }}
    />
    </motion.div>
    
    {!shouldReduceMotion && (
    <div className="absolute inset-0 z-0 opacity-10">
    <div className="grid grid-cols-12 h-full">
        {Array.from({ length: 12 }).map((_, i) => (
        <div key={i} className="border-r border-gray-400 h-full"></div>
        ))}
    </div>
    <div className="grid grid-rows-12 w-full absolute top-0 bottom-0">
        {Array.from({ length: 12 }).map((_, i) => (
        <div key={i} className="border-b border-gray-400 w-full"></div>
        ))}
    </div>
    </div>
    )}

    <div className="container mx-auto px-4 md:px-6 z-10">
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
        <motion.div
        initial="hidden"
        animate="visible"
        variants={contentVariants}
        >
        <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold mb-6">
            <motion.span 
            className="text-blue-600 dark:text-blue-400 inline-block"
            animate={shouldReduceMotion ? {} : { 
                color: ['#3b82f6', '#6366f1', '#4f46e5', '#3b82f6'],
            }}
            transition={shouldReduceMotion ? {} : { 
                duration: 8,
                repeat: Infinity,
                ease: "easeInOut"
            }}
            >
            {profile.name}
            </motion.span>
            <br />
            {profile.role}
        </h1>
        <p className="text-xl text-gray-600 dark:text-gray-300 mb-8 max-w-lg">
            {profile.bio.short}
        </p>
        <div className="flex flex-wrap gap-4">
            <Button
            onClick={() => setIsModalOpen(true)}
            >
            Предложить сотрудничество
            </Button>
        </div>
        </motion.div>

        <div className="relative mx-auto lg:mx-0 w-full max-w-md aspect-square">
        <div className="relative w-full h-full">
            <motion.div
            className="relative w-full h-full z-10"
            initial="hidden"
            animate="visible"
            variants={avatarVariants}
            >
            <motion.div
                className="absolute -inset-3 rounded-full border-2 border-blue-400 dark:border-blue-600 opacity-70"
                animate={shouldReduceMotion ? {} : { 
                rotate: 360
                }}
                transition={shouldReduceMotion ? {} : { 
                duration: 20,
                repeat: Infinity,
                ease: "linear"
                }}
                style={{ willChange: shouldReduceMotion ? 'auto' : 'transform' }}
            />
            
            <div className="relative w-full h-full rounded-full overflow-hidden border-4 border-blue-600 dark:border-blue-400 shadow-2xl">
                <Image
                src={profile.avatar}
                alt={profile.name}
                fill
                sizes="(max-width: 768px) 100vw, 400px"
                className="object-cover"
                priority
                />
            </div>
            </motion.div>
            
            {techIcons.map((icon, i) => {
            const angle = (i * (360 / techIcons.length)) % 360;
            const radius = shouldReduceMotion ? 200 : 360;
            
            // Round values to avoid precision differences
            const getPosition = (a: number) => ({
              x: Math.round(radius * Math.cos(a * Math.PI / 180) * 100) / 100,
              y: Math.round(radius * Math.sin(a * Math.PI / 180) * 100) / 100,
            });
                
                return (
                  <motion.div
                    key={icon.name}
                    className="absolute bg-white dark:bg-gray-800 p-2 rounded-full shadow-lg z-20"
                    style={{
                      width: icon.size,
                      height: icon.size,
                      top: '50%',
                      left: '50%',
                      transform: `translate(-50%, -50%)`,
                    }}
                    initial={{
                      ...getPosition(angle),
                      opacity: 0,
                      scale: shouldReduceMotion ? 1 : 0.5,
                    }}
                    animate={shouldReduceMotion ? {
                      opacity: 1,
                      scale: 1,
                      ...getPosition(angle),
                    } : {
                      opacity: 1,
                      scale: 1,
                      x: [
                        getPosition(angle).x,
                        getPosition(angle + 120).x,
                        getPosition(angle + 240).x,
                        getPosition(angle + 360).x,
                      ],
                      y: [
                        getPosition(angle).y,
                        getPosition(angle + 120).y,
                        getPosition(angle + 240).y,
                        getPosition(angle + 360).y,
                      ],
                    }}
                    transition={shouldReduceMotion ? {
                      duration: 0.5,
                      delay: i * 0.1,
                    } : {
                      duration: 20,
                      times: [0, 0.33, 0.66, 1],
                      repeat: Infinity,
                      delay: i * 0.2 + 0.8,
                      ease: "linear"
                    }}
                  >
                    <motion.div
                      animate={shouldReduceMotion ? {} : {
                        rotate: [-5, 5, -5],
                        scale: [0.95, 1.05, 0.95],
                      }}
                      transition={shouldReduceMotion ? {} : {
                        duration: 2,
                        repeat: Infinity,
                        repeatType: "mirror",
                        delay: i * 0.3,
                      }}
                    >
                      <Image
                        src={icon.src}
                        alt={icon.name}
                        width={icon.size - 4}
                        height={icon.size - 4}
                      />
                    </motion.div>
                  </motion.div>
                );
              })}
            </div>
          </div>
        </div>
      </div>

      {!shouldReduceMotion && Array.from({ length: 15 }).map((_, i) => {
        const size = (i % 8) + 10;
        const opacity = 0.2 + (i % 3) * 0.1;
        const startX = (i * 7) % 100;
        const startY = (i * 11) % 100;
        const duration = 12 + (i % 6) * 1.5;
        const xOffset = ((i % 5) - 2) * 20;
        
        return (
          <motion.div
            key={i}
            className="absolute rounded-full bg-blue-400 dark:bg-blue-500"
            style={{
              width: `${size}px`,
              height: `${size}px`,
              left: `${startX}%`,
              top: `${startY}%`,
              opacity: opacity,
              zIndex: 5,
              willChange: 'transform',
            }}
            animate={{
              y: [0, -300],
              x: [0, xOffset],
            }}
            transition={{
              duration: duration,
              repeat: Infinity,
              repeatType: "loop",
              ease: "linear",
              delay: i * 0.2,
            }}
          />
        );
      })}

      <CollaborationModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
      />
    </section>
  );
};
