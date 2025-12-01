export interface User {
  id: number;
  email: string;
  name: string;
  is_admin: boolean;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: number;
  name: string;
  user_id?: number | null;
  is_default: boolean;
  created_at: string;
  updated_at: string;
}

export interface Expense {
  id: number;
  name: string;
  category_id: number;
  category?: Category;
  user_id: number;
  unit: number;
  per_unit_cost: number;
  total: number;
  expense_date: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface AnalyticsData {
  total_expenses: number;
  expense_count: number;
  by_category: CategoryExpense[];
  daily_expenses: DailyExpense[];
  average_daily_spend: number;
  date_range: {
    start: string;
    end: string;
  };
}

export interface CategoryExpense {
  category_id: number;
  category_name: string;
  total: number;
  count: number;
}

export interface DailyExpense {
  date: string;
  total: number;
}

export interface DateRange {
  start: string;
  end: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
  total?: number;
  limit?: number;
  offset?: number;
}
