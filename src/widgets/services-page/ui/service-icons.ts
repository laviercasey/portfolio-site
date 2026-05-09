import { Bot, Code2, Database, Layers, Workflow } from 'lucide-react';
import type { LucideIcon } from 'lucide-react';
import type { ServiceIconKey } from '@/entities/service';

export const SERVICE_ICON_MAP: Record<ServiceIconKey, LucideIcon> = {
  bot: Bot,
  code: Code2,
  layers: Layers,
  database: Database,
  workflow: Workflow,
};
