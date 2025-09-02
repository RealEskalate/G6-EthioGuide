import { z } from 'zod';

export const newPasswordSchema = z
  .object({
    password: z
      .string()
      .min(8, { message: 'Password must be at least 8 characters' })
      .refine(
        (val) =>
          /[A-Z]/.test(val) && // uppercase
          /[a-z]/.test(val) && // lowercase
          /[0-9]/.test(val) && // number
          /[^A-Za-z0-9]/.test(val), // special character
        {
          message:
            'Password must contain at least one uppercase, one lowercase, one number, and one special character',
        }
      ),
    confirmPassword: z
      .string()
      // .min(8, { message: 'Confirm password must be at least 8 characters' }),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword'],
  });

export type NewPasswordFormData = z.infer<typeof newPasswordSchema>;
