'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuthStore } from '@/store/authStore';
import api from '@/lib/api';
import toast from 'react-hot-toast';
import { AnalyticsData, DateRange } from '@/types';
import { format } from 'date-fns';
import { PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8', '#82CA9D'];

export default function AnalyticsPage() {
  const router = useRouter();
  const { isAuthenticated, loadFromStorage } = useAuthStore();
  const [dateRange, setDateRange] = useState<DateRange | null>(null);
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [analytics, setAnalytics] = useState<AnalyticsData | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    loadFromStorage();
  }, [loadFromStorage]);

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/auth/login');
      return;
    }
    fetchDateRange();
  }, [isAuthenticated, router]);

  const fetchDateRange = async () => {
    try {
      const response = await api.get('/expenses/date-range');
      const range = response.data.data;
      setDateRange(range);
      
      // Set default dates
      setStartDate(format(new Date(range.start), 'yyyy-MM-dd'));
      setEndDate(format(new Date(range.end), 'yyyy-MM-dd'));
    } catch (error: any) {
      if (error.response?.status === 404) {
        toast.error('No expenses found. Add some expenses first!');
      } else {
        toast.error('Failed to load date range');
      }
    }
  };

  const fetchAnalytics = async () => {
    if (!startDate || !endDate) {
      toast.error('Please select start and end dates');
      return;
    }

    setIsLoading(true);
    try {
      const response = await api.get('/expenses/analytics', {
        params: {
          start_date: startDate,
          end_date: endDate,
        },
      });
      setAnalytics(response.data.data);
    } catch (error: any) {
      toast.error(error.response?.data?.error || 'Failed to load analytics');
    } finally {
      setIsLoading(false);
    }
  };

  if (!isAuthenticated) {
    return null;
  }

  const minDate = dateRange ? format(new Date(dateRange.start), 'yyyy-MM-dd') : '';
  const maxDate = format(new Date(), 'yyyy-MM-dd');

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <Link href="/dashboard" className="text-xl font-bold text-gray-900">
                MinFlow
              </Link>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          {/* Date Range Selector */}
          <div className="bg-white shadow rounded-lg p-6 mb-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">Select Date Range</h3>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
              <div>
                <label htmlFor="start-date" className="block text-sm font-medium text-gray-700">
                  Start Date
                </label>
                <input
                  type="date"
                  id="start-date"
                  min={minDate}
                  max={endDate || maxDate}
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm text-gray-900"
                  value={startDate}
                  onChange={(e) => setStartDate(e.target.value)}
                />
              </div>
              <div>
                <label htmlFor="end-date" className="block text-sm font-medium text-gray-700">
                  End Date
                </label>
                <input
                  type="date"
                  id="end-date"
                  min={startDate || minDate}
                  max={maxDate}
                  className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm text-gray-900"
                  value={endDate}
                  onChange={(e) => setEndDate(e.target.value)}
                />
              </div>
              <div className="flex items-end">
                <button
                  onClick={fetchAnalytics}
                  disabled={isLoading || !startDate || !endDate}
                  className="w-full px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                >
                  {isLoading ? 'Loading...' : 'Generate Analytics'}
                </button>
              </div>
            </div>
          </div>

          {/* Analytics Results */}
          {analytics && (
            <>
              {/* Summary Cards */}
              <div className="grid grid-cols-1 gap-6 sm:grid-cols-3 mb-6">
                <div className="bg-white shadow rounded-lg p-6">
                  <h4 className="text-sm font-medium text-gray-500">Total Expenses</h4>
                  <p className="mt-2 text-3xl font-semibold text-gray-900">
                    ₹{analytics.total_expenses.toFixed(2)}
                  </p>
                </div>
                <div className="bg-white shadow rounded-lg p-6">
                  <h4 className="text-sm font-medium text-gray-500">Number of Expenses</h4>
                  <p className="mt-2 text-3xl font-semibold text-gray-900">
                    {analytics.expense_count}
                  </p>
                </div>
                <div className="bg-white shadow rounded-lg p-6">
                  <h4 className="text-sm font-medium text-gray-500">Average Daily Spend</h4>
                  <p className="mt-2 text-3xl font-semibold text-gray-900">
                    ₹{analytics.average_daily_spend.toFixed(2)}
                  </p>
                </div>
              </div>

              {/* Charts */}
              <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
                {/* Category Breakdown */}
                <div className="bg-white shadow rounded-lg p-6">
                  <h3 className="text-lg font-medium text-gray-900 mb-4">Expenses by Category</h3>
                  {analytics.by_category.length > 0 ? (
                    <ResponsiveContainer width="100%" height={300}>
                      <PieChart>
                        <Pie
                          data={analytics.by_category as any}
                          dataKey="total"
                          nameKey="category_name"
                          cx="50%"
                          cy="50%"
                          outerRadius={100}
                          label
                        >
                          {analytics.by_category.map((entry, index) => (
                            <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                          ))}
                        </Pie>
                        <Tooltip />
                      </PieChart>
                    </ResponsiveContainer>
                  ) : (
                    <p className="text-gray-500 text-center py-8">No data available</p>
                  )}
                </div>

                {/* Daily Expenses */}
                <div className="bg-white shadow rounded-lg p-6">
                  <h3 className="text-lg font-medium text-gray-900 mb-4">Daily Expenses</h3>
                  {analytics.daily_expenses.length > 0 ? (
                    <ResponsiveContainer width="100%" height={300}>
                      <BarChart data={analytics.daily_expenses as any}>
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="date" />
                        <YAxis />
                        <Tooltip />
                        <Legend />
                        <Bar dataKey="total" fill="#8884d8" />
                      </BarChart>
                    </ResponsiveContainer>
                  ) : (
                    <p className="text-gray-500 text-center py-8">No data available</p>
                  )}
                </div>
              </div>

              {/* Category Table */}
              <div className="mt-6 bg-white shadow rounded-lg">
                <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
                  <h3 className="text-lg leading-6 font-medium text-gray-900">
                    Category Breakdown
                  </h3>
                </div>
                <div className="overflow-x-auto">
                  <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Category
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Count
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Total Amount
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                          Percentage
                        </th>
                      </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                      {analytics.by_category.map((category) => (
                        <tr key={category.category_id}>
                          <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                            {category.category_name}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {category.count}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            ₹{category.total.toFixed(2)}
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                            {((category.total / analytics.total_expenses) * 100).toFixed(1)}%
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            </>
          )}
        </div>
      </main>
    </div>
  );
}
