import { Github, Linkedin, Send, Mail, Twitter, Youtube, Globe } from 'lucide-react';
import {
  DiscordIcon,
  FacebookIcon,
  HabrIcon,
  InstagramIcon,
  VkIcon,
  WhatsAppIcon,
} from '@/shared/ui';

export const socialIconMap: Record<string, React.ComponentType<{ className?: string }>> = {
  github:    Github,
  linkedin:  Linkedin,
  habr:      HabrIcon,
  send:      Send,
  mail:      Mail,
  twitter:   Twitter,
  youtube:   Youtube,
  globe:     Globe,
  vk:        VkIcon,
  instagram: InstagramIcon,
  whatsapp:  WhatsAppIcon,
  discord:   DiscordIcon,
  facebook:  FacebookIcon,
};

