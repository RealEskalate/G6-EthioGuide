import { z } from 'zod';

export const registerSchema = z
  .object({
    fullName: z.string().min(2, 'Full name must be at least 2 characters').max(100, 'Full name must be less than 100 characters'),
    username: z
      .string()
      .min(3, 'Username must be at least 3 characters')
      .max(20, 'Username must be less than 20 characters')
      .regex(/^[a-zA-Z0-9_]+$/, 'Username can only contain letters, numbers, and underscores'),
    email: z.string().email('Invalid email address'),
    phoneNumber: z
      .string()
      .optional()
      .refine(
        (value) => !value || /^\+251[79]\d{6}(?:\d{2})?$/.test(value),
        {
          message: 'Phone number must be +251 followed by 7 or 9 digits starting with 7 or 9',
          path: ['phoneNumber'],
        }
      )
      .refine(
        (value) => !value || value.length >= 10,
        {
          message: 'Phone number must be at least 7 digits',
          path: ['phoneNumber'],
        }
      )
      .refine(
        (value) => !value || /^[79]/.test(value.slice(4)),
        {
          message: 'Phone number must start with 7 or 9',
          path: ['phoneNumber'],
        }
      ),
    password: z.string().min(8, 'Password must be at least 8 characters'),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  });

export type RegisterFormData = z.infer<typeof registerSchema>;