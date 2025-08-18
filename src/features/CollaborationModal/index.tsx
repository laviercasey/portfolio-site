"use client";

import { useState } from 'react';
import { motion } from 'framer-motion';
import { Modal } from '@/shared/ui/Modal';
import { Button } from '@/shared/ui/Button';
import { profile } from '@/entities/Profile';

interface CollaborationModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const CollaborationModal = ({
  isOpen,
  onClose,
}: CollaborationModalProps) => {
  const handleTelegramRedirect = () => {
    window.open(profile.social.telegram, '_blank');
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Предложить сотрудничество">
      <div className="text-center">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
        >
          <p className="text-gray-600 dark:text-gray-400 mb-6">
            Я открыт для новых проектов и возможностей сотрудничества. Давайте обсудим ваш проект в Telegram!
          </p>

          <Button onClick={handleTelegramRedirect} className="w-full">
            Перейти в Telegram
          </Button>
          
          <p className="mt-4 text-sm text-gray-500 dark:text-gray-400">
            Или напишите мне на email: {profile.email}
          </p>
        </motion.div>
      </div>
    </Modal>
  );
};
